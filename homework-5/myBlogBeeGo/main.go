/*
 * HomeWork-5: Start BeeGo
 * Created on 25.09.19 23:01
 * Copyright (c) 2019 - Eugene Klimov
 */

package main

import (
	"GoBasic/homework-5/myBlogBeeGo/models"
	_ "GoBasic/homework-5/myBlogBeeGo/routers"
	"database/sql"
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

const (
	DBNAME = "blog"
	DSN    = "/" + DBNAME + "?charset=utf8&interpolateParams=true"
)

func init() {
	var err error
	models.DB, err = sql.Open("mysql", myCnf("client")+DSN)
	if err != nil {
		log.Fatalln("Can't open DB:", err)
	}
	models.DB.SetMaxOpenConns(25)
	if err = models.DB.Ping(); err != nil {
		log.Fatalln("Can't ping DB:", err)
	}
}

func main() {
	defer models.DB.Close()
	beego.Run()
}
