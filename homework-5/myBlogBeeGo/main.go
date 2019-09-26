/*
 * HomeWork-5: Start BeeGo
 * Created on 25.09.19 23:01
 * Copyright (c) 2019 - Eugene Klimov
 */

package main

import (
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
	models.Lg.Info("Connected to BeeGo DB: %s", dbName)
	defer models.Lg.Close()

	beego.Run()
}
