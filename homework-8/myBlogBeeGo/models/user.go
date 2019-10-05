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
	"gopkg.in/mgo.v2/bson"
)

// User base struct
type User struct {
	Name string `json:"uname"`
	Pass string `json:"upass"`
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
	user := User{}
	err := d.Collection.FindOne(d.ctx, bson.M{"name": d.User.Name}).Decode(&user)
	if err == nil {
		return fmt.Errorf("user %s already exists", d.User.Name)
	}
	_, err = d.Collection.InsertOne(d.ctx, d.User)
	if err != nil {
		return fmt.Errorf("error create user: %v", err)
	}
	return nil
}

// GetUser read user data.
func (d *DBUsers) GetUser() error {
	user := User{}
	err := d.Collection.FindOne(d.ctx, bson.M{"name": d.User.Name, "pass": d.User.Pass}).Decode(&user)
	if err != nil {
		return fmt.Errorf("user %s not found", d.User.Name)
	}
	return nil
}