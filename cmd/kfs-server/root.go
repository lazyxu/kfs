package main

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/lazyxu/kfs/dao"
	"github.com/lazyxu/kfs/db/cgosqlite"
	"github.com/lazyxu/kfs/db/mysql"
	"github.com/lazyxu/kfs/db/sqlc"
	storage "github.com/lazyxu/kfs/storage/local"

	"github.com/gorilla/websocket"

	"github.com/lazyxu/kfs/rpc/server"

	"github.com/lazyxu/kfs/core"

	"github.com/spf13/cobra"
)

const (
	SocketServerStr   = "socket-server"
	WebServerStr      = "web-server"
	DatabaseTypeStr   = "database-type"
	DataSourceNameStr = "data-source-name"
	StorageTypeStr    = "storage-type"
	StorageDirStr     = "storage-dir"
	ThumbnailDirStr   = "thumbnail-dir"
	TransCodeDirStr   = "transcode-dir"
)

func init() {
	rootCmd.PersistentFlags().String(WebServerStr, "1123", "web server port")
	rootCmd.PersistentFlags().String(SocketServerStr, "1124", "socket server port")
	rootCmd.PersistentFlags().String(DatabaseTypeStr, "sqlite", "sqlite/mysql")
	rootCmd.PersistentFlags().String(DataSourceNameStr, "kfs.db", "data source name")
	rootCmd.PersistentFlags().String(StorageTypeStr, "1", "storage type, [0, 5]")
	rootCmd.PersistentFlags().String(StorageDirStr, "tmp", "storage dir path")
	rootCmd.PersistentFlags().String(ThumbnailDirStr, "thumbnail", "thumbnail dir path")
	rootCmd.PersistentFlags().String(TransCodeDirStr, "transcode", "transcode dir path")
}

//go:embed build/*
var build embed.FS

func getStorageByType(typ string) (func(string) (dao.Storage, error), error) {
	switch typ {
	case "0":
		return storage.NewStorage0, nil
	case "1":
		return storage.NewStorage1, nil
	case "2":
		return storage.NewStorage2, nil
	case "3":
		return storage.NewStorage3, nil
	case "4":
		return storage.NewStorage4, nil
	case "5":
		return storage.NewStorage5, nil
	}
	return nil, fmt.Errorf("no such storage type: %s", typ)
}

func getDatabaseByType(typ string) (func(string) (dao.Database, error), error) {
	switch typ {
	case "sqlite":
		return cgosqlite.New, nil
	case "mysql":
		return mysql.New, nil
	}
	return nil, fmt.Errorf("no such databse type: %s", typ)
}

var kfsCore *core.KFS

var rootCmd = &cobra.Command{
	Use:   "kfs-server",
	Short: "Kfs is file system used to backup files.",
	Args:  cobra.RangeArgs(0, 0),
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		defer func() {
			if err != nil {
				cmd.PrintErrln(err)
				os.Exit(1)
			}
		}()

		storageType := cmd.Flag(StorageTypeStr).Value.String()
		storageDir := cmd.Flag(StorageDirStr).Value.String()

		databaseType := cmd.Flag(DatabaseTypeStr).Value.String()
		dataSourceName := cmd.Flag(DataSourceNameStr).Value.String()

		var q *sqlc.Queries
		if dataSourceName != "" {
			// 自动建库逻辑，仅对 MySQL 有效
			if databaseType == "mysql" {
				// 提取不带 db 的 DSN
				// 自动获取数据库名
				noDBDSN := dataSourceName
				dbName := ""
				if idx := strings.LastIndex(noDBDSN, "/"); idx != -1 {
					rest := noDBDSN[idx+1:]
					qidx := strings.Index(rest, "?")
					if qidx != -1 {
						dbName = rest[:qidx]
						noDBDSN = noDBDSN[:idx+1] + rest[qidx:]
					} else {
						dbName = rest
						noDBDSN = noDBDSN[:idx+1]
					}
				}
				if dbName == "" {
					fmt.Fprintln(os.Stderr, "excect database name from data source name failed:", dataSourceName)
					os.Exit(1)
				}
				initDB, err := sql.Open("mysql", noDBDSN)
				if err == nil {
					initDB.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s`", dbName))
					initDB.Close()
				}
			}
			db, err := sql.Open("mysql", dataSourceName)
			if err != nil {
				fmt.Fprintln(os.Stderr, "connect error:", err)
				os.Exit(1)
			}
			defer db.Close()
			// 执行 db\sqlc\schema.sql 中的建表语句
			schemaSQL, err := os.ReadFile("db/sqlc/schema.sql")
			if err != nil {
				fmt.Fprintln(os.Stderr, "read schema.sql error:", err)
				os.Exit(1)
			}
			sqls := strings.Split(string(schemaSQL), ";")
			for _, sqlStmt := range sqls {
				sqlStmt = strings.TrimSpace(sqlStmt)
				if sqlStmt == "" {
					continue
				}
				_, err = db.Exec(sqlStmt)
				if err != nil {
					fmt.Fprintln(os.Stderr, "exec schema.sql error:", err, "SQL:", sqlStmt)
					os.Exit(1)
				}
			}
			if err != nil {
				fmt.Fprintln(os.Stderr, "exec schema.sql error:", err)
				os.Exit(1)
			}
			q = sqlc.New(db)
			// 插入用户
			_, err = q.InsertUser(context.Background(), sqlc.InsertUserParams{
				Email:    "test@user.com",
				Password: "password123",
			})
			if err != nil && !strings.Contains(err.Error(), "Duplicate entry") {
				fmt.Fprintln(os.Stderr, "InsertUser error:", err)
				os.Exit(1)
			}

			// 查询用户
			user, err := q.GetUserByEmail(context.Background(), "test@user.com")
			if err != nil {
				fmt.Fprintln(os.Stderr, "GetUserByName error:", err)
			}
			fmt.Println("User:", user)
			userID := user.ID
			// 插入文件
			res, err := q.InsertFile(context.Background(), sqlc.InsertFileParams{
				Size:   12345,
				Md5:    "d41d8cd98f00b204e980067ae2aefaef",
				Sha256: "e3b0c44298fc6cdcd5eaab0e63b8b8f6c7d4fb5c0f6e8d1c4f8b2a3d9f8b8f6c",
			})
			if err != nil {
				fmt.Fprintln(os.Stderr, "InsertFile error:", err)
				os.Exit(1)
			}
			fileID, err := res.LastInsertId()
			if err != nil {
				fmt.Fprintln(os.Stderr, "LastInsertId error:", err)
				os.Exit(1)
			}
			fmt.Println("Inserted file ID:", fileID)

			// 查询文件
			logicFiles, err := q.ListLogicFileByUserId(context.Background(), userID)
			if err != nil {
				fmt.Fprintln(os.Stderr, "GetLogicFileByPath error:", err)
			}
			fmt.Printf("logicFiles: %+v\n", logicFiles)

			// 更新文件
			err = q.SetLogicFileByPath(context.Background(), 1, []string{"test", "path"}, "ttt", sql.NullInt64{Int64: fileID, Valid: true}, 0644, 12345, 1622547800, 1622547800, 1622547800, 1622547800, true)
			if err != nil {
				fmt.Fprintln(os.Stderr, "SetLogicFileByPath error:", err)
				// 查询文件
				logicFiles, err = q.ListLogicFileByUserId(context.Background(), userID)
				if err != nil {
					fmt.Fprintln(os.Stderr, "GetLogicFileByPath error:", err)
				}
				fmt.Printf("logicFiles: %+v\n", logicFiles)
				os.Exit(1)
			}
			fmt.Println("Logic file updated")

			// 查询文件
			logicFiles, err = q.ListLogicFileByUserId(context.Background(), userID)
			if err != nil {
				fmt.Fprintln(os.Stderr, "GetLogicFileByPath error:", err)
			}
			fmt.Printf("logicFiles: %+v\n", logicFiles)
		}
		webPortString := cmd.Flag(WebServerStr).Value.String()
		webPort, err := strconv.Atoi(webPortString)
		if err != nil {
			return
		}
		if webPort != 0 && webPort < 1024 || webPort > 65535 {
			err = errors.New("webPort should be between 1024 and 65535, actual " + webPortString)
			return
		}

		socketPortString := cmd.Flag(SocketServerStr).Value.String()
		socketPort, err := strconv.Atoi(socketPortString)
		if err != nil {
			return
		}
		if socketPort != 0 && socketPort < 1024 || socketPort > 65535 {
			err = errors.New("socketPort should be between 1024 and 65535, actual " + socketPortString)
			return
		}

		newStorage, err := getStorageByType(storageType)
		if err != nil {
			return
		}
		newDatabase, err := getDatabaseByType(databaseType)
		if err != nil {
			return
		}

		thumbnailDir := cmd.Flag(ThumbnailDirStr).Value.String()
		transCodeDir := cmd.Flag(TransCodeDirStr).Value.String()
		kfsCore, err = core.New3(q, dao.DatabaseNewFunc(dataSourceName, newDatabase), dao.StorageNewFunc(storageDir, newStorage), thumbnailDir, transCodeDir)
		if err != nil {
			return
		}
		defer kfsCore.Close()

		startAllCloudSync()
		hostname, err := os.Hostname()
		if err != nil {
			return
		}
		println(hostname)
		dir, err := os.UserConfigDir()
		if err != nil {
			return
		}
		println(dir)
		go func() {
			// socket
			lis, err := net.Listen("tcp", "0.0.0.0:"+socketPortString)
			if err != nil {
				panic(err)
			}
			println("KFS socket server listening at:", lis.Addr().String())
			err = server.Socket(lis, kfsCore)
			if err != nil {
				panic(err)
			}
		}()
		// web
		http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
			wsHandler(w, r, kfsCore)
		})
		webServer(webPortString)
	},
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

func wsHandler(w http.ResponseWriter, r *http.Request, kfsCore *core.KFS) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	server.Process(kfsCore, ToAddrReadWriteCloser(c))
}
