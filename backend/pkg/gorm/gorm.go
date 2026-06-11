package gorm

import (
	"GO1/global"
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
		log.Fatalf("mysql connection failed: %v", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour * 4)

	sqlDir := global.Config.SqlFile.Path
	files, err := filepath.Glob(filepath.Join(sqlDir, "*.sql"))
	if err != nil {
		global.Logger.Fatalf("SQL file glob failed: %v", err)
	}
	sort.Strings(files)

	for _, file := range files {
		sqlBytes, err := os.ReadFile(file)
		if err != nil {
			global.Logger.Fatalf("SQL file read failed: %s, error: %v", file, err)
		}

		statements := splitSQLStatements(string(sqlBytes))
		if len(statements) == 0 {
			continue
		}

		for index, statement := range statements {
			if err := db.Exec(statement).Error; err != nil {
				global.Logger.Fatalf("SQL execution failed: %s statement %d, error: %v", file, index+1, err)
			}
		}

		global.Logger.Infof("Successfully executed SQL file: %s (%d statements)", file, len(statements))
	}

	global.DB = db
}

func splitSQLStatements(sql string) []string {
	var statements []string
	var builder strings.Builder
	inSingleQuote := false
	inDoubleQuote := false
	inBacktick := false
	inLineComment := false
	inBlockComment := false

	for i := 0; i < len(sql); i++ {
		ch := sql[i]
		var next byte
		if i+1 < len(sql) {
			next = sql[i+1]
		}

		if inLineComment {
			builder.WriteByte(ch)
			if ch == '\n' {
				inLineComment = false
			}
			continue
		}

		if inBlockComment {
			builder.WriteByte(ch)
			if ch == '*' && next == '/' {
				builder.WriteByte(next)
				i++
				inBlockComment = false
			}
			continue
		}

		if !inSingleQuote && !inDoubleQuote && !inBacktick {
			switch {
			case ch == '-' && next == '-':
				builder.WriteByte(ch)
				builder.WriteByte(next)
				i++
				inLineComment = true
				continue
			case ch == '#':
				builder.WriteByte(ch)
				inLineComment = true
				continue
			case ch == '/' && next == '*':
				builder.WriteByte(ch)
				builder.WriteByte(next)
				i++
				inBlockComment = true
				continue
			}
		}

		switch ch {
		case '\'':
			builder.WriteByte(ch)
			if !inDoubleQuote && !inBacktick && !isEscaped(sql, i) {
				if inSingleQuote && next == '\'' {
					builder.WriteByte(next)
					i++
					continue
				}
				inSingleQuote = !inSingleQuote
			}
			continue
		case '"':
			if !inSingleQuote && !inBacktick && !isEscaped(sql, i) {
				inDoubleQuote = !inDoubleQuote
			}
		case '`':
			if !inSingleQuote && !inDoubleQuote {
				inBacktick = !inBacktick
			}
		case ';':
			if !inSingleQuote && !inDoubleQuote && !inBacktick {
				appendSQLStatement(&statements, builder.String())
				builder.Reset()
				continue
			}
		}

		builder.WriteByte(ch)
	}

	appendSQLStatement(&statements, builder.String())
	return statements
}

func appendSQLStatement(statements *[]string, statement string) {
	statement = strings.TrimSpace(statement)
	if statement != "" {
		*statements = append(*statements, statement)
	}
}

func isEscaped(sql string, pos int) bool {
	backslashes := 0
	for i := pos - 1; i >= 0 && sql[i] == '\\'; i-- {
		backslashes++
	}
	return backslashes%2 == 1
}
