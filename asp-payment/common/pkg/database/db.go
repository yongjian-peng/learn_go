package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	DB    *gorm.DB
	SqlDB *sql.DB
)

func Init(user, password, host, port, dbName string, maxIdleConns, maxOpenConns int) {
	//日志
	dbLogger := logger.New(
		log.New(&lumberjack.Logger{
			Filename:   "./logs/db.log", // 日志文件路径
			MaxSize:    10,              // 最大M
			MaxBackups: 5,               // 最多保留多少个备份
			MaxAge:     24,              // days
			Compress:   false,           // 是否压缩 disabled by default
		}, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // Log level
			Colorful:      false,       // 禁用彩色打印
		},
	)
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, dbName)
	if DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true, //禁用 创建、更新、删除 默认事务
		Logger:                 dbLogger,
		PrepareStmt:            true, //缓存预编译语句
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
		},
	}); err != nil {
		panic("数据库连接异常")
	}

	if SqlDB, err = DB.DB(); err != nil {
		panic("获取sqlDb异常")
	}

	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	SqlDB.SetMaxIdleConns(maxIdleConns)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	SqlDB.SetMaxOpenConns(maxOpenConns)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	SqlDB.SetConnMaxLifetime(time.Hour)

}
