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

// DB & Logger (temporary)
var DB *sql.DB
var Lg *logs.BeeLogger

const (
	TABLENAME   = "posts"
	GETALLPOSTS = "SELECT id, title, summary, body, DATE_FORMAT(updated, '%d.%m.%Y %H:%i') FROM " + TABLENAME + " WHERE deleted IS NULL ORDER BY id DESC"
	GETONEPOST  = "SELECT id, title, summary, body, DATE_FORMAT(updated, '%d.%m.%Y %H:%i') FROM " + TABLENAME + " WHERE deleted IS NULL AND id = ?"
	//DELETEPOST   = "DELETE FROM " + TABLENAME + " WHERE id = ?"
	DELETEPOST = "UPDATE " + TABLENAME + " SET deleted = ? WHERE id = ?"
	INSERTPOST = "INSERT INTO " + TABLENAME + " (title, summary, body) VALUES(?, ?, ?)"
	UPDATEPOST = "UPDATE " + TABLENAME + " SET title = ?, summary = ?, body = ? WHERE ID = ?"
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
	DB    *sql.DB
	Posts []Post
	Lg    *logs.BeeLogger
	Error
}

// NewPosts creates new DBPosts with DB link
func NewPosts() *DBPosts {
	return &DBPosts{
		DB:    DB,
		Lg:    Lg,
		Error: Error{Lg: Lg},
	}
}

// GetPosts gets one or all posts.
func (p *DBPosts) GetPosts(id string) error {
	var rows *sql.Rows
	var err error
	if id != "" {
		rows, err = p.DB.Query(GETONEPOST, id)
	} else {
		rows, err = p.DB.Query(GETALLPOSTS)
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

// create one post.
func (p *DBPosts) createPost(post *Post, db *sql.DB) error {
	_, err := db.Exec(INSERTPOST, post.Title, post.Summary, post.Body)
	return err
}

// DeletePost deletes one post.
func (p *DBPosts) DeletePost(id string) error {
	delTime := time.Now().Format("2006-01-02 15:04:05")
	_, err := p.DB.Exec(DELETEPOST, delTime, id)
	return err
}

// update one post.
func (p *DBPosts) updatePost(post *Post, db *sql.DB) error {
	_, err := db.Exec(UPDATEPOST, post.Title, post.Summary, post.Body, post.ID)
	return err
}
