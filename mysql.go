package utils

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	sql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	maxId   = 5
	maxOpen = 10
)

func MysqlConnect(address, username, password, database string) (*gorm.DB, error) {
	// 对密码进行URL编码，处理特殊字符（如@、:、/等）
	log.Println(fmt.Sprintf(
		"connect mysql - %s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, address, database,
	))

	// 创建 GORM 日志记录器
	conn, err := gorm.Open(
		sql.Open(fmt.Sprintf(
			"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			username, password, address, database,
		)),
		&gorm.Config{
			Logger: logger.New(
				log.StandardLogger(), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
				logger.Config{
					SlowThreshold: time.Second, // 慢 SQL 阈值
					LogLevel:      logger.Info, // 日志级别
					// IgnoreRecordNotFoundError: true,          // 忽略ErrRecordNotFound（记录未找到）错误
					Colorful: false, // 禁用彩色打印
				},
			),
		},
	)
	if err != nil {
		return nil, err
	}

	db, err := conn.DB()
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(maxId)             // SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	db.SetMaxOpenConns(maxOpen)           // SetMaxOpenConns sets the maximum number of open connections to the database.
	db.SetConnMaxLifetime(24 * time.Hour) // SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	db.SetConnMaxIdleTime(time.Hour)

	return conn, nil
}
