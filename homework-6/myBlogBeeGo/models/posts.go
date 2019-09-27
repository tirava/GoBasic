/*
 * HomeWork-6: Mongo in BeeGo
 * Created on 28.09.19 22:16
 * Copyright (c) 2019 - Eugene Klimov
 */

package models

import (
	"context"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"html/template"
	"time"
)

// MDB & Logger & ORM are globals (it is normal for BeeGo)
var (
	Lg  *logs.BeeLogger
	ORM orm.Ormer
	MDB *mongo.Client
)

// Post is the base post type.
type Post struct {
	ID string `json:"-" bson:"_id,omitempty"`
	//ID      int           `orm:"column(id);pk;auto"`
	Title   string        `json:"title"`
	Date    time.Time     `json:"-" orm:"column(updated_at);auto_now;type(datetime)"`
	Summary string        `json:"summary"`
	Body    template.HTML `json:"body"`
	Created time.Time     `json:"-" orm:"column(created_at);auto_now_add;type(datetime)"`
	Deleted time.Time     `json:"-" orm:"column(deleted_at);type(datetime)"`
}

//DBPosts is type dbPosts map[string]Post
type DBPosts struct {
	//DB     *sql.DB
	MDB    *mongo.Client
	DBName string
	//ORM    orm.Ormer
	Posts []Post
	Lg    *logs.BeeLogger
	Error
}

// NewPosts creates new DBPosts with DB link
func NewPosts() *DBPosts {
	return &DBPosts{
		Lg:    Lg,
		Error: Error{Lg: Lg},
		//ORM:    ORM,
		MDB:    MDB,
		DBName: beego.AppConfig.String("DBNAME"),
	}
}

// TableName returns name for post's table (need for ORM & Mongo).
func (Post) TableName() string {
	table := beego.AppConfig.String("TABLENAME")
	if table == "" {
		return "posts"
	}
	return table
}

// Date2Norm normalizes date to local format for view in browsers.
func (p *Post) Date2Norm() string {
	dt := beego.AppConfig.String("DATETIME")
	if dt == "" {
		dt = "02.01.2006 15:04:05"
	}
	return p.Date.Format(dt)
}

//// GetPosts gets one or all posts.
//func (d *DBPosts) GetPosts(id string) error {
//	if id == "" { // all posts
//		qs := d.ORM.QueryTable(&Post{})
//		n, err := qs.Filter("deleted_at__isnull", true).OrderBy("-updated_at").All(&d.Posts)
//		if err != nil {
//			return fmt.Errorf("error in query all posts: %v", err)
//		}
//		if n == 0 {
//			d.Lg.Error("no one posts found")
//			return nil
//		}
//	} else { // one post
//		idn, err := strconv.Atoi(id)
//		if err != nil {
//			return fmt.Errorf("error while converting post ID: %s", id)
//		}
//		post := Post{ID: idn}
//		if err := d.ORM.Read(&post); err != nil {
//			return fmt.Errorf("post not found: %s", id)
//		}
//		d.Posts = append(d.Posts, post)
//	}
//	return nil
//}

// GetPosts gets one or all posts.
func (d *DBPosts) GetPosts(id string) error {
	post := Post{}
	c := d.MDB.Database(d.DBName).Collection(post.TableName())
	if id == "" { // all posts
		cur, err := c.Find(context.TODO(), bson.D{})
		if err != nil {
			return fmt.Errorf("error in query all posts: %v", err)
		}
		err = cur.All(context.TODO(), &d.Posts)
	} else { // one post
		err := c.FindOne(context.TODO(), bson.D{{"_id", id}}).Decode(&post)
		if err != nil {
			return fmt.Errorf("post not found: %s", id)
		}
		d.Posts = append(d.Posts, post)
	}
	return nil
}

//// CreatePost creates post.
//func (d *DBPosts) CreatePost() error {
//	n, err := d.ORM.Insert(&d.Posts[0])
//	if n == 0 {
//		return fmt.Errorf("post not created")
//	}
//	return err
//}

// CreatePost creates post.
func (d *DBPosts) CreatePost() error {
	d.Posts[0].Date = time.Now()
	d.Posts[0].Created = time.Now()
	d.Posts[0].Deleted = time.Unix(0, 0)
	c := d.MDB.Database(d.DBName).Collection(d.Posts[0].TableName())
	_, err := c.InsertOne(context.TODO(), d.Posts[0])
	return err
}

// DeletePost deletes one post.
func (d *DBPosts) DeletePost(id string) error {
	//qs := d.ORM.QueryTable(&Post{})
	//n, err := qs.Filter("id", id).Update(orm.Params{"deleted_at": time.Now().Local()})
	//if n == 0 {
	//	return fmt.Errorf("post not found: %s", id)
	//}
	//return err
	return nil
}

// UpdatePost updates post.
func (d *DBPosts) UpdatePost() error {
	//n, err := d.ORM.Update(&d.Posts[0])
	//if n == 0 {
	//	return fmt.Errorf("post not found: %d", d.Posts[0].ID)
	//}
	//return err
	return nil
}
