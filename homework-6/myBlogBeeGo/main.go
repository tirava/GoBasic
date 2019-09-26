/*
 * HomeWork-5: Start BeeGo
 * Created on 25.09.19 23:01
 * Copyright (c) 2019 - Eugene Klimov
 */

package main

import (
	"context"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	// connect to Mongo
	mdb, err := mongo.NewClient(options.Client().ApplyURI("mongodb://Klim.Go:27017"))
	if err != nil {
		log.Fatalln("Can't create MongoDB client:", err)
	}
	models.MDB = mdb
	if err = models.MDB.Connect(context.TODO()); err != nil {
		log.Fatalln("Can't connect to MongoDB server:", err)
	}
	if err = models.MDB.Ping(context.TODO(), nil); err != nil {
		log.Fatalln("Can't ping MongoDB server:", err)
	}
	defer models.MDB.Disconnect(context.TODO())

	// set logger
	dbName := beego.AppConfig.String("DBNAME")
	models.Lg = logs.NewLogger(10)
	models.Lg.SetPrefix(fmt.Sprintf("[%s]", dbName))
	models.Lg.Info("Connected to BeeGo DB: %s", dbName)
	models.Lg.Info("Connected to MongoDB")
	defer models.Lg.Close()

	beego.Run()
}
