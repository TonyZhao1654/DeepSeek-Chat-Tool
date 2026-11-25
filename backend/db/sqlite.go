package db

/*
	该文件不进行错误包装
	错误包装在本包的db.go
*/

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// newSQLite 打开 SQLite 连接并完成初始化
func newSQLite(path string) (*DB, error) {
	// 把相对路径转成绝对路径，避免相对路径的歧义问题
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	// 取出数据库文件所在目录
	dir := filepath.Dir(absolutePath)
	err = os.MkdirAll(dir, 0o755) // 如果目录不存在就递归创建
	if err != nil {
		return nil, err
	}

	// 打开sqlite数据库文件
	sqlDB, err := sql.Open("sqlite", absolutePath)
	if err != nil {
		return nil, err
	}

	//同时保证 PRAGMA 对所有操作生效
	sqlDB.SetMaxOpenConns(1) // 这就是个本地单用户工具，单连接足够

	// 测试连接是否可用
	err = sqlDB.Ping()
	if err != nil {
		sqlDB.Close() // 关闭连接池
		return nil, err
	}

	// 构造DB结构体
	database := &DB{DB: sqlDB}

	// 配置
	err = database.setup()
	if err != nil {
		sqlDB.Close() // 关闭连接池
		return nil, err
	}

	return database, nil
}

func (d *DB) setup() error {
	pragmas := []string{
		"PRAGMA busy_timeout = 5000", // 忙等待超时时间
		"PRAGMA journal_mode = WAL",  // 日志模式: WAL
		"PRAGMA foreign_keys = ON",   // 外键约束
	}
	/*
		WAL(Write-Ahead Logging，预写式日志)
		读和写可以并发，大幅改善sqlite读写性能
	*/

	// 逐个执行PRAGMA语句
	for _, pragma := range pragmas {
		_, err := d.Exec(pragma)
		if err != nil {
			return err
		}
	}

	return nil
}
