package main

import (
	"context"

	"github.com/leeexeo/kon/log"
	"github.com/leeexeo/kon/orm"
)

var logConf = &log.Config{
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

var dbConf = &orm.Config{
	User:        "root",
	Password:    "12345678",
	Host:        "localhost",
	Port:        "3306",
	Schema:      "user",
	MaxIdleConn: 20,
	MaxOpenConn: 5,
	Charset:     "utf8",
	Engine:      "InnoDB",
	Collate:     "utf8_bin",
}

type User struct {
	Name string
	Sex  string
	Age  int
}

func main() {
	log.SetupGlobal(logConf)
	ctx := context.Background()
	log.Debug(ctx, "test debug log", "a", "b")
	log.Info(ctx, "text info log", "c", "d", "e")

	orm.SetupGlobal(dbConf, logConf, User{})
	db := orm.GetDb()
	db.Create(&User{
		Name: "lixiang",
		Sex:  "man",
		Age:  26,
	})
}
