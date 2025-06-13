package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"hobbyhub-server/models"
)

var DB *gorm.DB

// CreateDatabase 创建数据库（如果不存在）
func CreateDatabase(conf *Config) error {
	switch conf.Database.Type {
	case "mysql":
		// 连接到MySQL服务器（不指定数据库）
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/",
			conf.Database.Username,
			conf.Database.Password,
			conf.Database.Host,
			conf.Database.Port)

		db, err := sql.Open("mysql", dsn)
		if err != nil {
			return fmt.Errorf("连接MySQL服务器失败: %v", err)
		}
		defer db.Close()

		// 创建数据库
		_, err = db.Exec("CREATE DATABASE IF NOT EXISTS " + conf.Database.Database)
		if err != nil {
			return fmt.Errorf("创建数据库失败: %v", err)
		}

		log.Printf("数据库 %s 已创建或已存在\n", conf.Database.Database)
	case "sqlite":
		// 对于SQLite，确保数据库目录存在
		dbDir := filepath.Dir(conf.Database.Database)
		if dbDir != "" && dbDir != "." {
			if err := os.MkdirAll(dbDir, 0755); err != nil {
				return fmt.Errorf("创建SQLite数据库目录失败: %v", err)
			}
		}
		log.Printf("SQLite数据库文件路径已准备: %s\n", conf.Database.Database)
	default:
		return fmt.Errorf("不支持的数据库类型: %s", conf.Database.Type)
	}

	return nil
}

func InitDatabase(conf *Config) error {
	// 首先尝试创建数据库
	if err := CreateDatabase(conf); err != nil {
		return fmt.Errorf("创建数据库失败: %v", err)
	}

	var err error

	switch conf.Database.Type {
	case "mysql":
		// 如果是 MySQL 数据库，使用 gorm 连接
		username := conf.Database.Username
		password := conf.Database.Password
		host := conf.Database.Host
		port := conf.Database.Port
		database := conf.Database.Database

		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			username, password, host, port, database, conf.Database.Charset)
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return fmt.Errorf("连接数据库失败: %v", err)
		}
		log.Printf("成功连接到MySQL数据库 %s (地址: %s:%d, 用户: %s)\n", database, host, port, username)
	case "sqlite":
		// 如果是 SQLite 数据库，使用 sqlite 驱动
		DB, err = gorm.Open(sqlite.Open(conf.Database.Database), &gorm.Config{})
		if err != nil {
			return fmt.Errorf("连接SQLite数据库失败: %v", err)
		}
		log.Printf("成功连接到SQLite数据库: %s\n", conf.Database.Database)
	default:
		return fmt.Errorf("不支持的数据库类型: %s", conf.Database.Type)
	}

	// 自动迁移所有模型
	err = DB.AutoMigrate(
		&models.User{},
		&models.Friend{},
		&models.File{},
		&models.Chat{},
		&models.Activity{},
		&models.ActivityMember{},
		&models.ActivityComment{},
		&models.Admin{},
	)

	if err != nil {
		return fmt.Errorf("自动迁移模型失败: %v", err)
	}
	log.Println("所有模型已成功迁移到数据库")

	return nil
}
