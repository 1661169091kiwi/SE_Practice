package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	mysql "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// Init 初始化数据库连接
func Init(dsn string) (*sql.DB, error) {
	var err error

	// 创建数据库连接
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// 设置连接池参数
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(25)
	DB.SetConnMaxLifetime(5 * time.Minute)

	// 测试连接
	if err = DB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	log.Println("Database connection established successfully")
	return DB, nil
}

// GetDB 获取数据库连接实例
func GetDB() *sql.DB {
	return DB
}

// Close 关闭数据库连接
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// HealthCheck 数据库健康检查
func HealthCheck() error {
	if DB == nil {
		return fmt.Errorf("database connection is nil")
	}
	return DB.Ping()
}

// QueryRow 查询单行数据
func QueryRow(query string, args ...interface{}) *sql.Row {
	if DB == nil {
		return nil
	}
	return DB.QueryRow(query, args...)
}

// Query 查询多行数据
func Query(query string, args ...interface{}) (*sql.Rows, error) {
	if DB == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	return DB.Query(query, args...)
}

// Exec 执行SQL语句（INSERT, UPDATE, DELETE）
func Exec(query string, args ...interface{}) (sql.Result, error) {
	if DB == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	return DB.Exec(query, args...)
}

// BeginTransaction 开始事务
func BeginTransaction() (*sql.Tx, error) {
	if DB == nil {
		return nil, fmt.Errorf("database connection is nil")
	}
	return DB.BeginTx(context.Background(), &sql.TxOptions{Isolation: sql.LevelReadCommitted})
}

// Insert 插入数据并返回自增ID
func Insert(query string, args ...interface{}) (int64, error) {
	result, err := Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// Update 更新数据并返回影响行数
func Update(query string, args ...interface{}) (int64, error) {
	result, err := Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Delete 删除数据并返回影响行数
func Delete(query string, args ...interface{}) (int64, error) {
	result, err := Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// IsDuplicateEntryError 检查是否为重复条目错误
func IsDuplicateEntryError(err error) bool {
	if err == nil {
		return false
	}
	var me *mysql.MySQLError
	if errors.As(err, &me) {
		return me.Number == 1062
	}
	s := err.Error()
	return strings.Contains(s, "1062") || strings.Contains(s, "Duplicate entry")
}

// IsForeignKeyError 检查是否为外键约束错误
func IsForeignKeyError(err error) bool {
	if err == nil {
		return false
	}
	var me *mysql.MySQLError
	if errors.As(err, &me) {
		return me.Number == 1452
	}
	s := err.Error()
	return strings.Contains(s, "1452") || strings.Contains(s, "foreign key constraint fails")
}

// TableExists 检查表是否存在
func TableExists(tableName string) (bool, error) {
	query := "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?"
	var count int
	err := QueryRow(query, tableName).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetTableRowCount 获取表的行数
func GetTableRowCount(tableName string) (int64, error) {
	query := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
	var count int64
	err := QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
