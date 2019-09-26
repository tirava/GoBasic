/*
 * HomeWork-5: Start BeeGo
 * Created on 25.09.19 23:06
 * Copyright (c) 2019 - Eugene Klimov
 */

package models

import (
	"database/sql"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"html/template"
	"strconv"
	"time"
)

// Logger & ORM are globals (it is normal for BeeGo)
var (
	Lg  *logs.BeeLogger
	ORM orm.Ormer
)

// Post is the base post type.
type Post struct {
	ID      int           `orm:"column(id);pk;auto"`
	Title   string        `json:"title"`
	Date    time.Time     `json:"-" orm:"column(updated_at);auto_now;type(datetime)"`
	Summary string        `json:"summary"`
	Body    template.HTML `json:"body"`
	Created time.Time     `json:"-" orm:"column(created_at);auto_now_add;type(datetime)"`
	Deleted time.Time     `json:"-" orm:"column(deleted_at);type(datetime)"`
}

//DBPosts is type dbPosts map[string]Post
type DBPosts struct {
	DB    *sql.DB
	ORM   orm.Ormer
	Posts []Post
	Lg    *logs.BeeLogger
	Error
}

// NewPosts creates new DBPosts with DB link
func NewPosts() *DBPosts {
	return &DBPosts{
		Lg:    Lg,
		Error: Error{Lg: Lg},
		ORM:   ORM,
	}
}

// TableName returns name for post's table (need for ORM).
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

// GetPosts gets one or all posts.
func (d *DBPosts) GetPosts(id string) error {
	if id == "" { // all posts
		qs := d.ORM.QueryTable(&Post{})
		n, err := qs.Filter("deleted_at__isnull", true).OrderBy("-updated_at").All(&d.Posts)
		if err != nil {
			return fmt.Errorf("error in query all posts: %v", err)
		}
		if n == 0 {
			d.Lg.Error("no one posts found")
			return nil
		}
	} else { // one post
		idn, err := strconv.Atoi(id)
		if err != nil {
			return fmt.Errorf("error while converting post ID: %s", id)
		}
		post := Post{ID: idn}
		if err := d.ORM.Read(&post); err != nil {
			return fmt.Errorf("post not found: %s", id)
		}
		d.Posts = append(d.Posts, post)
	}
	return nil
}

// CreatePost creates post.
func (d *DBPosts) CreatePost() error {
	n, err := d.ORM.Insert(&d.Posts[0])
	if n == 0 {
		return fmt.Errorf("post not created")
	}
	return err
}

// DeletePost deletes one post.
func (d *DBPosts) DeletePost(id string) error {
	qs := d.ORM.QueryTable(&Post{})
	n, err := qs.Filter("id", id).Update(orm.Params{"deleted_at": time.Now().Local()})
	if n == 0 {
		return fmt.Errorf("post not found: %s", id)
	}
	return err
}

// UpdatePost updates post.
func (d *DBPosts) UpdatePost() error {
	n, err := d.ORM.Update(&d.Posts[0])
	if n == 0 {
		return fmt.Errorf("post not found: %d", d.Posts[0].ID)
	}
	return err
}
