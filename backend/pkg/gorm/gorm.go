package gorm

import (
	"GO1/global"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitGorm() {

	dsn := global.Config.Mysql.DSN()

	var mysqlLogger logger.Interface
	if global.Config.Mysql.LogLevel == "debug" {
		mysqlLogger = logger.Default.LogMode(logger.Info)
	} else {
		mysqlLogger = logger.Default.LogMode(logger.Error)
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: mysqlLogger,
	})
	if err != nil {
		log.Fatal(fmt.Sprintf("[%s] mysql连接失败！", dsn))
		return
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour * 4)

	files := []string{
		"database/sql/calendar_submissions.sql",
		"database/sql/competitions.sql",
	}

	var allSQL string
	for _, file := range files {
		sqlBytes, err := os.ReadFile(file)
		if err != nil {
			global.Logger.Error("SQL数据读取失败！文件：" + file)
			return
		}

		allSQL += string(sqlBytes) + "\n"
	}

	err = db.Exec(allSQL).Error
	if err != nil {
		global.Logger.Error("数据库创建失败！")
		return
	}

	global.DB = db
}
