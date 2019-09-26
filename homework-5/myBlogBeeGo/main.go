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
	"github.com/astaxie/beego/orm"
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
	cnf := beego.AppConfig.String("MYCNFPROFILE")
	if cnf == "" {
		cnf = "client"
	}

	// connect to DB
	models.DB, err = sql.Open("mysql", myCnf(cnf)+"/"+dbName+dsn)
	if err != nil {
		log.Fatalln("Can't open DB:", err)
	}
	defer models.DB.Close()

	if err = models.DB.Ping(); err != nil {
		log.Fatalln("Can't ping DB:", err)
	}
	models.DB.SetMaxOpenConns(25)

	// BeeGo ORM
	err = orm.RegisterDataBase("default", "mysql", myCnf("client")+"/"+dbName+dsn)
	if err != nil {
		log.Fatalln("Can't open BeeGo DB:", err)
	}
	orm.RegisterModel(new(models.Post))
	models.ORM = orm.NewOrm()

	// set logger
	models.Lg = logs.NewLogger(10)
	models.Lg.SetPrefix("[" + dbName + "]")
	models.Lg.Info("Connected to DB: %s", dbName)
	models.Lg.Info("Connected to BeeGo DB: %s", dbName)
	defer models.Lg.Close()

	orm.Debug = true
	beego.Run()
}
