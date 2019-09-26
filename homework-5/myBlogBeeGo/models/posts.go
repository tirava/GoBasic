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

//// Constatnts.
//const (
//	DELDATETMPL = "2006-01-02 15:04:05"
//)

// DB, Logger, ORM are globals (it is normal for BeeGo)
var (
	DB  *sql.DB
	Lg  *logs.BeeLogger
	ORM orm.Ormer
)

// Post is the base post type.
type Post struct {
	ID      int           `orm:"column(id);pk;auto"`
	Title   string        `json:"title"`
	Date    time.Time     `json:"-" orm:"column(updated_at);auto_now"`
	Summary string        `json:"summary"`
	Body    template.HTML `json:"body"`
	Created time.Time     `json:"-" orm:"column(created_at);auto_now"`
	Deleted time.Time     `json:"-" orm:"column(deleted_at)"`
}

//DBPosts is type dbPosts map[string]Post
type DBPosts struct {
	DB *sql.DB
	DBQueries
	ORM   orm.Ormer
	Posts []Post
	Lg    *logs.BeeLogger
	Error
}

// NewPosts creates new DBPosts with DB link
func NewPosts() *DBPosts {
	return &DBPosts{
		DB:        DB,
		Lg:        Lg,
		Error:     Error{Lg: Lg},
		DBQueries: *NewDBQueries(),
		ORM:       ORM,
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

// GetPosts gets one or all posts.
func (p *DBPosts) GetPosts(id string) error {
	ids, _ := strconv.Atoi(id)
	post := Post{ID: ids}
	if id == "" { // all posts
		qs := p.ORM.QueryTable(&post)
		n, err := qs.Filter("deleted_at__isnull", true).OrderBy("-updated_at").All(&p.Posts)
		if err != nil {
			return fmt.Errorf("error in query all posts: %v", err)
		}
		if n == 0 {
			p.Lg.Error("no one posts found")
			return nil
		}
	} else { // one post
		if err := p.ORM.Read(&post); err != nil {
			return fmt.Errorf("post not found: %s", id)
		}
		p.Posts = append(p.Posts, post)
	}
	return nil
}

// CreatePost creates post.
func (p *DBPosts) CreatePost() error {
	_, err := p.ORM.Insert(&p.Posts[0])
	return err
}

//// DeletePost deletes one post.
//func (p *DBPosts) DeletePost(id string) error {
//	delTime := time.Now().Format(DELDATETMPL)
//	_, err := p.DB.Exec(p.QDeletePost, delTime, id)
//	return err
//}

// DeletePost deletes one post.
func (p *DBPosts) DeletePost(id string) error {
	qs := p.ORM.QueryTable(&Post{})
	_, err := qs.Filter("id", id).Update(orm.Params{"deleted_at": time.Now().Local()})
	return err
}

// UpdatePost updates post.
func (p *DBPosts) UpdatePost() error {
	_, err := p.DB.Exec(p.QUpdatePost, p.Posts[0].Title, p.Posts[0].Summary, p.Posts[0].Body, p.Posts[0].ID)
	return err
}
