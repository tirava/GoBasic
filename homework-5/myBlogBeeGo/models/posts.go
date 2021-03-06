/*
 * HomeWork-5: Start BeeGo
 * Created on 25.09.19 23:06
 * Copyright (c) 2019 - Eugene Klimov
 */

package models

import (
	"database/sql"
	"fmt"
	"github.com/astaxie/beego/logs"
	"html/template"
	"time"
)

// Constatnts.
const (
	DELDATETMPL = "2006-01-02 15:04:05"
)

// DB & Logger global (it is normal for BeeGo)
var (
	DB *sql.DB
	Lg *logs.BeeLogger
)

// Post is the base post type.
type Post struct {
	ID      string
	Title   string        `json:"title"`
	Date    string        `json:"date"`
	Summary string        `json:"summary"`
	Body    template.HTML `json:"body"`
}

//DBPosts is type dbPosts map[string]Post
type DBPosts struct {
	DB *sql.DB
	DBQueries
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
	}
}

// GetPosts gets one or all posts.
func (p *DBPosts) GetPosts(id string) error {
	var rows *sql.Rows
	var err error
	if id != "" {
		rows, err = p.DB.Query(p.QGetOnePost, id)
	} else {
		rows, err = p.DB.Query(p.QGetAllPosts)
	}
	defer rows.Close()
	if err != nil {
		return fmt.Errorf("error in db.query: %v", err)
	}
	for rows.Next() {
		post := Post{}
		err = rows.Scan(&post.ID, &post.Title, &post.Summary, &post.Body, &post.Date)
		if err != nil {
			return fmt.Errorf("error in row.scan: %v", err)
		}
		p.Posts = append(p.Posts, post)
	}
	if len(p.Posts) == 0 {
		return fmt.Errorf("post not found: %s", id)
	}
	return nil
}

// CreatePost creates post.
func (p *DBPosts) CreatePost() error {
	_, err := p.DB.Exec(p.QInsertPost, p.Posts[0].Title, p.Posts[0].Summary, p.Posts[0].Body)
	return err
}

// DeletePost deletes one post.
func (p *DBPosts) DeletePost(id string) error {
	delTime := time.Now().Format(DELDATETMPL)
	_, err := p.DB.Exec(p.QDeletePost, delTime, id)
	return err
}

// UpdatePost updates post.
func (p *DBPosts) UpdatePost() error {
	_, err := p.DB.Exec(p.QUpdatePost, p.Posts[0].Title, p.Posts[0].Summary, p.Posts[0].Body, p.Posts[0].ID)
	return err
}
