package gorm

import (
	"GO1/global"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
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

	sqlDir := global.Config.SqlFile.Path

	files, err := filepath.Glob(filepath.Join(sqlDir, "*.sql"))
	if err != nil {
		global.Logger.Error("SQL文件匹配失败: " + err.Error())
		return
	}
	// 保证每次执行顺序一致
	sort.Strings(files)

	for _, file := range files {
		sqlBytes, err := os.ReadFile(file)
		if err != nil {
			global.Logger.Error("SQL文件读取失败: " + file + ", 错误: " + err.Error())
			return
		}

		sql := strings.TrimSpace(string(sqlBytes))
		if sql == "" {
			continue
		}

		if err := db.Exec(sql).Error; err != nil {
			global.Logger.Error("执行SQL失败: " + file + ", 错误: " + err.Error())
			return
		}

		global.Logger.Info("成功执行SQL文件: " + file)
	}

	global.DB = db
}
