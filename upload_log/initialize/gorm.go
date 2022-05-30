package initialize

// package gorm

import (
	"log"
	"upload_log/global"

	// "bop_notify/initialize/internal"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

//@author: SliverHorn
//@function: Gorm
//@description: 初始化数据库并产生数据库全局变量
//@return: *gorm.DB
func Gorm() *gorm.DB {
	// fmt.Println(global.GLOBAL_CONFIG.System.DbType)
	// switch global.GLOBAL_CONFIG.System.DbType {
	// case "mysql":
	// 	return GormMysql()
	// default:
	// 	return GormMysql()
	// }
	// return GormMysql()
	return Conn()
}

func Conn() *gorm.DB {
	m := global.Config.Mysql
	// fmt.Println("hhh", m)
	if m.Dbname == "" {
		return nil
	}
	dsn := m.Username + ":" + m.Password + "@tcp(" + m.Path + ")/" + m.Dbname + "?" + m.Config
	// fmt.Println(dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			// TablePrefix:   "bop_", // 表名前缀，`Article` 的表名应该是 `it_articles`
			SingularTable: true, // 使用单数表名，启用该选项，此时，`Article` 的表名应该是 `it_article`
		},
	})

	// fmt.Println(db, err)
	if err != nil {
		log.Panic("mysql err:", err)
		return nil
	} else {
		return db
	}
}

// NamingStrategy: schema.NamingStrategy{SingularTable: true}
// 例： db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{NamingStrategy: schema.NamingStrategy{SingularTable: true}})

//
//@author: SliverHorn
//@function: GormMysql
//@description: 初始化Mysql数据库
//@return: *gorm.DB

func GormMysql() *gorm.DB {
	m := global.Config.Mysql
	fmt.Println("mysql的配置", m)
	if m.Dbname == "" {
		return nil
	}
	dsn := m.Username + ":" + m.Password + "@tcp(" + m.Path + ")/" + m.Dbname + "?" + m.Config
	mysqlConfig := mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         191,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
		// NamingStrategy: schema.NamingStrategy{
		// 	// TablePrefix:   "bop_", // 表名前缀，`Article` 的表名应该是 `it_articles`
		// 	SingularTable: true, // 使用单数表名，启用该选项，此时，`Article` 的表名应该是 `it_article`
		// },
	}
	// fmt.Println("HHHHH")
	if db, err := gorm.Open(mysql.New(mysqlConfig), gormConfig(m.LogMode)); err != nil {
		fmt.Println("mysql异常")
		//global.GVA_LOG.Error("MySQL启动异常", zap.Any("err", err))
		//os.Exit(0)
		//return nil
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		// fmt.Println(db)
		return db
	}
}

//@author: SliverHorn
//@function: gormConfig
//@description: 根据配置决定是否开启日志
//@param: mod bool
//@return: *gorm.Config

func gormConfig(mod bool) *gorm.Config {
	var config = &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true}
	switch global.Config.Mysql.LogZap {
	// case "silent", "Silent":
	// 	config.Logger = internal.Default.LogMode(logger.Silent)
	// case "error", "Error":
	// 	config.Logger = internal.Default.LogMode(logger.Error)
	// case "warn", "Warn":
	// 	config.Logger = internal.Default.LogMode(logger.Warn)
	// case "info", "Info":
	// 	config.Logger = internal.Default.LogMode(logger.Info)
	// case "zap", "Zap":
	// 	config.Logger = internal.Default.LogMode(logger.Info)
	// default:
	// 	if mod {
	// 		config.Logger = internal.Default.LogMode(logger.Info)
	// 		break
	// 	}
	// 	config.Logger = internal.Default.LogMode(logger.Silent)
	}
	return config
}
