package config

import (
	"database/sql"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"hobbyhub-server/models"
)

var DB *gorm.DB

// CreateDatabase 创建数据库（如果不存在）
func CreateDatabase(conf *Config) error {
	if conf.Database.Type != "mysql" {
		return fmt.Errorf("不支持的数据库类型: %s", conf.Database.Type)
	}

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

	fmt.Printf("数据库 %s 已创建或已存在", conf.Database.Database)
	return nil
}

func InitDatabase(conf *Config) error {
	// 首先尝试创建数据库
	if err := CreateDatabase(conf); err != nil {

		return fmt.Errorf("创建数据库失败: %v", err)
	}

	username := conf.Database.Username
	password := conf.Database.Password
	host := conf.Database.Host
	port := conf.Database.Port
	database := conf.Database.Database

	if conf.Database.Type == "mysql" {
		// 如果是 MySQL 数据库，使用 gorm 连接
		var err error
		dsn := username + ":" + password + "@tcp(" + host + ":" + fmt.Sprint(port) + ")/" + database + "?charset=" + conf.Database.Charset + "&parseTime=True&loc=Local"
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {

			return fmt.Errorf("连接数据库失败: %v\n", err)
		}
		fmt.Printf("成功连接到MySQL数据库 %s (地址: %s:%d, 用户: %s)", database, host, port, username)

		// 自动迁移所有模型
		err = DB.AutoMigrate(
			&models.User{},
			&models.Friend{},
			&models.File{},
			&models.Chat{},
			&models.Activity{},
			&models.ActivityMember{},
			&models.ActivityComment{},
		)

		if err != nil {

			return fmt.Errorf("自动迁移模型失败: %v", err)
		}
		fmt.Println("所有模型已成功迁移到数据库")

		return nil
	} else {
		return fmt.Errorf("不支持的数据库类型: %s", conf.Database.Type)
	}
}
