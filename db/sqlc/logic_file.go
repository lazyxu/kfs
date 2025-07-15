package sqlc

import (
	"context"
	"database/sql"
	"fmt"
	"os"
)

func (q *Queries) GetLogicFileByPath(ctx context.Context, userID int64, dirpath []string) (*LogicFile, error) {
	var parentID sql.NullInt64
	var lf *LogicFile

	for _, name := range dirpath {
		params := GetLogicFileParams{
			UserID:   userID,
			ParentID: parentID,
			Name:     name,
		}
		logicFile, err := q.GetLogicFile(ctx, params)
		if err != nil {
			return nil, err
		}
		lf = &logicFile
		parentID.Scan(lf.ID)
	}

	if lf == nil {
		return nil, sql.ErrNoRows
	}

	return lf, nil
}

// 设置某路径的LogicFile，如果某一层不存在则创建（mkdir控制），只处理_logic_file表
func (q *Queries) SetLogicFileByPath(ctx context.Context, userID int64, dirpath []string, name string, fileID sql.NullInt64, mode, size, createTime, modifyTime, changeTime, accessTime int64, mkdir bool) error {
	var parentID sql.NullInt64
	var lf *LogicFile

	for _, name := range dirpath {
		params := GetLogicFileParams{
			UserID:   userID,
			ParentID: parentID,
			Column3:  parentID,
			Name:     name,
		}
		logicFile, err := q.GetLogicFile(ctx, params)
		if err != nil {
			if err == sql.ErrNoRows && mkdir {
				createParams := InsertLogicFileParams{
					ParentID:   parentID,
					UserID:     userID,
					Name:       name,
					FileID:     sql.NullInt64{Valid: false},
					Mode:       int64(os.ModeDir | 0755),
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
					return err
				}
				fmt.Printf("Inserted LogicFile with ParentID: %v, UserID: %d, Name: %s\n", parentID, userID, name)
				fmt.Printf("GetLogicFile: %+v\n", params)
				logicFile, err = q.GetLogicFile(ctx, params)
				if err != nil {
					fmt.Fprintf(os.Stderr, "GetLogicFile 2 error: %v", err)
					return err
				}
			} else {
				fmt.Fprintf(os.Stderr, "GetLogicFile 1 error: %v", err)
				return err
			}
		}
		lf = &logicFile
		parentID.Scan(lf.ID)
	}

	if lf == nil {
		fmt.Fprintf(os.Stderr, "sql.ErrNoRows error")
		return sql.ErrNoRows
	}

	// 可选：更新最后一层的文件属性
	updateParams := UpdateLogicFileParams{
		ParentID:   lf.ParentID,
		UserID:     lf.UserID,
		Name:       name,
		FileID:     fileID,
		Mode:       mode,
		Size:       size,
		CreateTime: createTime,
		ModifyTime: modifyTime,
		ChangeTime: changeTime,
		AccessTime: accessTime,
		ID:         lf.ID,
	}
	if err := q.UpdateLogicFile(ctx, updateParams); err != nil {
		fmt.Fprintf(os.Stderr, "UpdateLogicFile error: %v", err)
		return err
	}

	return nil
}

func (q *Queries) ListDir(ctx context.Context, userID int64, dirpath []string) (*LogicFile, []*LogicFile, error) {
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
	children, err := q.ListLogicFileByParentId(ctx, parentID)
	if err != nil {
		return lf, nil, err
	}

	// 转换为 []*LogicFile
	childPtrs := make([]*LogicFile, 0, len(children))
	for i := range children {
		childPtrs = append(childPtrs, &children[i])
	}

	return lf, childPtrs, nil
}
