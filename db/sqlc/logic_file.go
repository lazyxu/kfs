package sqlc

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

func (q *Queries) GetLogicFileByPath(ctx context.Context, userID int64, dirpath []string) (*LogicFile, *File, error) {
	var parentID sql.NullInt64
	var lf *LogicFile

	for _, name := range dirpath {
		params := GetLogicFileParams{
			UserID:   userID,
			ParentID: parentID,
			Name:     name,
		}
		log.Printf("Getting LogicFile for UserID: %d, ParentID: %v, Name: %s\n", userID, parentID, name)
		logicFile, err := q.GetLogicFile(ctx, params)
		if err != nil {
			return nil, nil, err
		}
		log.Printf("logicFile: %+v\n", logicFile)
		lf = &logicFile
		parentID.Scan(lf.ID)
	}

	log.Printf("lf: %+v\n", lf)
	if lf == nil {
		return nil, nil, sql.ErrNoRows
	}
	if !lf.FileID.Valid {
		return lf, nil, nil // directory
	}
	f, err := q.GetFileByID(ctx, lf.FileID.Int64)
	if err != nil {
		return nil, nil, err
	}
	log.Printf("f: %+v\n", f)
	return lf, &f, nil
}

const DefaultFileMode int64 = 0o644
const DefaultDirMode = int64(os.ModeDir | 0755)

// 设置某路径的LogicFile，如果某一层不存在则创建（mkdir控制），只处理_logic_file表
func (q *Queries) SetLogicFileByPath(ctx context.Context, userID int64, dirpath []string, name string, fileID sql.NullInt64,
	mode, size int64, createTime, modifyTime, changeTime, accessTime time.Time, mkdir bool) (int64, error) {
	var parentID sql.NullInt64

	for _, name := range dirpath {
		params := GetLogicFileParams{
			UserID:   userID,
			ParentID: parentID,
			Name:     name,
		}
		logicFile, err := q.GetLogicFile(ctx, params)
		if err != nil {
			if err != sql.ErrNoRows {
				fmt.Fprintf(os.Stderr, "GetLogicFile 1 error: %v", err)
				return 0, err
			} else if mkdir {
				createParams := InsertLogicFileParams{
					UserID:     userID,
					ParentID:   parentID,
					Name:       name,
					FileID:     sql.NullInt64{Valid: false},
					Mode:       DefaultDirMode,
					Size:       0,
					CreateTime: createTime,
					ModifyTime: modifyTime,
					ChangeTime: changeTime,
					AccessTime: accessTime,
				}
				fmt.Printf("Inserting LogicFile: %+v\n", createParams)
				_, err = q.InsertLogicFile(ctx, createParams)
				if err != nil {
					fmt.Fprintf(os.Stderr, "InsertLogicFile error: %v", err)
					return 0, err
				}
				fmt.Printf("Inserted LogicFile with ParentID: %v, UserID: %d, Name: %s\n", parentID, userID, name)
				fmt.Printf("GetLogicFile: %+v\n", params)
				logicFile, err = q.GetLogicFile(ctx, params)
				if err != nil {
					fmt.Fprintf(os.Stderr, "GetLogicFile 2 error: %v", err)
					return 0, err
				}
			}
		}
		parentID.Scan(logicFile.ID)
	}

	// 可选：更新最后一层的文件属性
	updateParams := UpsertLogicFileParams{
		UserID:     userID,
		ParentID:   parentID,
		Name:       name,
		FileID:     fileID,
		Mode:       mode,
		Size:       size,
		CreateTime: createTime,
		ModifyTime: modifyTime,
		ChangeTime: changeTime,
		AccessTime: accessTime,
	}
	log.Printf("Upserting LogicFile: %+v\n", updateParams)
	r, err := q.UpsertLogicFile(ctx, updateParams)
	if err != nil {
		fmt.Fprintf(os.Stderr, "UpdateLogicFile error: %v", err)
		return 0, err
	}

	return r.LastInsertId()
}

func (q *Queries) ListDir(ctx context.Context, userID int64, dirpath []string) (*LogicFile, []*ListLogicFilesByParentIdRow, error) {
	var parentID sql.NullInt64
	var lf *LogicFile

	// 逐层查找目录
	for _, name := range dirpath {
		params := GetLogicFileParams{
			UserID:   userID,
			ParentID: parentID,
			Name:     name,
		}
		logicFile, err := q.GetLogicFile(ctx, params)
		if err != nil {
			return nil, nil, err
		}
		lf = &logicFile
		parentID.Scan(lf.ID)
	}

	// 列出该目录下的所有子项
	children, err := q.ListLogicFilesByParentId(ctx, ListLogicFilesByParentIdParams{UserID: userID, ParentID: parentID})
	if err != nil {
		return lf, nil, err
	}

	// 转换为 []*LogicFile
	childPtrs := make([]*ListLogicFilesByParentIdRow, 0, len(children))
	for i := range children {
		childPtrs = append(childPtrs, &children[i])
	}

	return lf, childPtrs, nil
}
