/*
 * HomeWork-8: Config, Logs and Auth
 * Created on 06.10.19 12:16
 * Copyright (c) 2019 - Eugene Klimov
 */

package models

import (
	"context"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"go.mongodb.org/mongo-driver/mongo"
)

// User base struct
type User struct {
	Name string
	Pass string
}

//DBUsers is the base type for posts
type DBUsers struct {
	Collection *mongo.Collection
	ctx        context.Context
	User       User
	Lg         *logs.BeeLogger
	Error
}

// TableName returns name for user's table (need for ORM & Mongo).
func (User) TableName() string {
	table := beego.AppConfig.String("USERSTABLE")
	if table == "" {
		return "users"
	}
	return table
}

// NewUser creates new User instance.
func NewUser() *DBUsers {
	dbName := beego.AppConfig.String("DBNAME")
	col := MDB.Database(dbName).Collection(User{}.TableName())
	return &DBUsers{
		Collection: col,
		ctx:        context.TODO(),
		Lg:         Lg,
		Error:      Error{Lg: Lg},
	}
}

// CreateUser creates user.
func (d *DBUsers) CreateUser() error {
	fmt.Println(d.User)
	//d.Posts[0].OID = primitive.NewObjectID() // or omitempty in Post
	//d.Posts[0].Date = time.Now()
	//d.Posts[0].Created = time.Now()
	//d.Posts[0].Deleted = time.Unix(0, 0)
	//_, err := d.Collection.InsertOne(d.ctx, d.Posts[0])
	//if err != nil {
	//	return fmt.Errorf("error insert one post: %v", err)
	//}
	return nil
}
