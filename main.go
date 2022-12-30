package main

import (
	"context"
	"fmt"

	"github.com/leeexeo/kon/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var config = &log.Config{
	Filename:                  "./test.log",
	MaxSize:                   1,
	MaxAge:                    1,
	MaxBackups:                1,
	LocalTime:                 false,
	Compress:                  false,
	CallerSkip:                3,
	Level:                     "info",
	Console:                   "stdout",
	GormLevel:                 "info",
	SqlSlowThreshold:          200,
	IgnoreRecordNotFoundError: false,
	IgnoreDuplicateError:      false,
}

type User struct {
	Name string
	Sex  string
	Age  int
}

func main() {
	log.SetupGlobal(config)
	ctx := context.Background()
	log.Debug(ctx, "test debug log", "a", "b")
	log.Info(ctx, "text info log", "c", "d", "e")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true", "root", "12345678", "127.0.0.1", 3306, "user", "utf8")
	log.Debug(ctx, dsn)
	gormLogger, _ := log.NewGormLogger(config)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		panic("failed to connect database")
	}
	err = db.AutoMigrate(User{})
	if err != nil {
		panic("failed to AutoMigrate")
	}
	db.Create(&User{
		Name: "lixiang",
		Sex:  "man",
		Age:  26,
	})
}
