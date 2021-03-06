/*
 * HomeWork-5: Start BeeGo
 * Created on 25.09.19 23:01
 * Copyright (c) 2019 - Eugene Klimov
 */

package main

import (
	"database/sql"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"myBlog/models"
	_ "myBlog/routers"
)

func main() {
	var err error
	dbName := beego.AppConfig.String("DBNAME")
	if dbName == "" {
		dbName = "blog"
	}
	dsn := beego.AppConfig.String("DSN")
	if dsn == "" {
		dsn = "?charset=utf8&interpolateParams=true"
	}

	// connect to DB
	models.DB, err = sql.Open("mysql", myCnf("client")+"/"+dbName+dsn)
	if err != nil {
		log.Fatalln("Can't open DB:", err)
	}
	defer models.DB.Close()

	if err = models.DB.Ping(); err != nil {
		log.Fatalln("Can't ping DB:", err)
	}
	models.DB.SetMaxOpenConns(25)

	// set logger
	models.Lg = logs.NewLogger(10)
	models.Lg.SetPrefix("[" + dbName + "]")
	models.Lg.Info("Connected to DB: %s", dbName)
	defer models.Lg.Close()

	beego.Run()
}
