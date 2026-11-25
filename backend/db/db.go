package db

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
)

// sqlFiles内嵌sql文件夹下的建表语句，按文件名顺序执行
//
//go:embed sql/*.sql
var sqlFiles embed.FS

type DB struct {
	*sql.DB // 数据库句柄
}

func NewDB(path string) (*DB, error) {
	sqlite, err := newSQLite(path)
	if err != nil {
		return nil, fmt.Errorf("db: new db: %w", err)
	}

	// 执行sql文件夹下的sql文件
	err = initSchema(sqlite)
	if err != nil {
		sqlite.Close() // 关闭连接池

		// initSchema包装好错误了，这里直接返回
		return nil, err
	}

	return sqlite, nil
}

func initSchema(d *DB) error {
	// 读取虚拟目录
	entries, err := fs.ReadDir(sqlFiles, "sql")
	if err != nil {
		return fmt.Errorf("db: read sql dir: %w", err)
	}

	for _, entry := range entries {
		// 过滤子目录
		if entry.IsDir() {
			continue // 遍历到子文件夹就跳过(忽略子文件夹)
		}

		// 拼接出文件在embed.FS中的完整相对路径
		name := "sql/" + entry.Name()

		// 读取为字节切片
		statement, err := sqlFiles.ReadFile(name)
		if err != nil {
			return fmt.Errorf("db: read %s: %w", name, err)
		}

		_, err = d.Exec(string(statement))
		if err != nil {
			return fmt.Errorf("db: exec %s: %w", name, err)
		}
	}

	return nil
}
