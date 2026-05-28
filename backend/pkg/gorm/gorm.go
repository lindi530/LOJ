package gorm

import (
	"GO1/global"
	"fmt"
	"log"
	"os"
	"strings"
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
		"database/sql/competition_players.sql",
		"database/sql/competition_problems.sql",
		"database/sql/competition_submissions.sql",
	}

	for _, file := range files {
		sqlBytes, err := os.ReadFile(file)
		if err != nil {
			global.Logger.Error("SQL文件读取失败: " + file)
			return
		}
		sql := string(sqlBytes)
		// 去除首尾空白，如果末尾有分号，MySQL 可以接受（也可以去掉）
		sql = strings.TrimSpace(sql)
		if sql == "" {
			continue
		}
		if err := db.Exec(sql).Error; err != nil {
			global.Logger.Error("执行SQL失败: " + file + ", 错误: " + err.Error())
			return
		}
		global.Logger.Info("成功执行SQL文件: " + file)
	}

	if err != nil {
		global.Logger.Error("数据库创建失败！")
		return
	}

	global.DB = db
}
