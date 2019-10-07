/*
 * HomeWork-7: Testing & Docs in BeeGo
 * Created on 28.09.19 22:18
 * Copyright (c) 2019 - Eugene Klimov
 */

package main

import (
	"context"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"myBlogBeeGo/conf"
	"myBlogBeeGo/models"
	_ "myBlogBeeGo/routers"
)

func main() {
	// connect to Mongo
	mdb, err := mongo.NewClient(options.Client().ApplyURI(conf.GetURI()))
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

	// set logger
	dbName := beego.AppConfig.String("DBNAME")
	models.Lg = logs.NewLogger(10)
	models.Lg.SetPrefix(fmt.Sprintf("[%s]", dbName))
	models.Lg.Info("Connected to MongoDB")
	logFile := beego.AppConfig.String("LOGFILE")
	logParam := fmt.Sprintf(`{"filename":"%s"}`, logFile)
	err = models.Lg.SetLogger(logs.AdapterConsole)
	err = models.Lg.SetLogger(logs.AdapterFile, logParam)
	if err != nil {
		log.Printf("can't set logger write to file: %s\n", logFile)
	}
	defer models.Lg.Close()

	beego.BConfig.Listen.Graceful = true
	beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	beego.Run()

	if err = models.MDB.Disconnect(context.TODO()); err != nil {
		log.Fatalln("error disconnect from MongoDB", err)
	}
}
