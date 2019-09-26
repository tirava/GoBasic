/*
 * HomeWork-5: Start BeeGo
 * Created on 25.09.19 23:01
 * Copyright (c) 2019 - Eugene Klimov
 */

package main

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"myBlog/conf"
	"myBlog/models"
	_ "myBlog/routers"
)

func main() {
	// BeeGo ORM
	err := orm.RegisterDataBase("default", "mysql", conf.GetDSN())
	if err != nil {
		log.Fatalln("Can't open BeeGo DB:", err)
	}
	orm.RegisterModel(new(models.Post))
	models.ORM = orm.NewOrm()

	// set logger
	dbName := beego.AppConfig.String("DBNAME")
	models.Lg = logs.NewLogger(10)
	models.Lg.SetPrefix(fmt.Sprintf("[%s]", dbName))
	models.Lg.Info("Connected to BeeGo DB: %s", dbName)
	defer models.Lg.Close()

	beego.Run()
}
