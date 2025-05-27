package config

import (
	"log"

	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"hobbyhub-server/models"
)

var DB *gorm.DB

func InitDatabase(conf *Config) {
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
			log.Fatalf("Failed to connect to database: %v", err)
		}
		log.Printf("Connected to MySQL database %s at %s:%d as user %s", database, host, port, username)

		/*
			tableNameList := []string{}
			rows, err := DB.Raw("SHOW TABLES").Rows()
			if err != nil {
				log.Fatalf("Failed to list tables: %v", err)
			}
			defer rows.Close()
			for rows.Next() {
				var tableName string
				if err := rows.Scan(&tableName); err != nil {
					log.Printf("Failed to scan table name: %v", err)
					continue
				}
				tableNameList = append(tableNameList, tableName)
			}
			log.Printf("Tables in database: %v", tableNameList)
		*/
		var user models.User
		if err := DB.AutoMigrate(&user); err != nil {
			log.Fatalf("Failed to auto migrate User model: %v", err)
		}
		return
	} else {
		log.Fatalf("Unsupported database type: %s", conf.Database.Type)
		return
	}
}
