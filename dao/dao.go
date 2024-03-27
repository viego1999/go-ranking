package dao

import (
	"time"

	"go-ranking/config"
	"go-ranking/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	Db  *gorm.DB
	err error
)

func init() {
	Db, err = gorm.Open(mysql.Open(config.Mysqldb), &gorm.Config{})
	if err != nil {
		logger.Error(map[string]interface{}{"mysql connect error": err.Error()})
	}
	if Db != nil {
		logger.Error(map[string]interface{}{"database error": Db.Error})
	}
	sqlDB, sqlErr := Db.DB()
	if sqlErr != nil {
		logger.Error(map[string]interface{}{"sql db error": sqlErr.Error()})
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
}
