/*
 * HomeWork-5: Start BeeGo
 * Created on 25.09.19 23:06
 * Copyright (c) 2019 - Eugene Klimov
 */

package models

import (
	"database/sql"
	"fmt"
	"html/template"
	"time"
)

// DB is global DB var (temporary)
var DB *sql.DB

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
}

// GetPosts gets one or all posts.
func (p DBPosts) GetPosts(id string) (DBPosts, error) {
	var rows *sql.Rows
	var err error
	var posts = DBPosts{}
	if id != "" {
		rows, err = p.DB.Query(GETONEPOST, id)
	} else {
		rows, err = p.DB.Query(GETALLPOSTS)
	}
	if err != nil {
		return posts, fmt.Errorf("error in db.query: %v", err)
	}
	for rows.Next() {
		post := Post{}
		err = rows.Scan(&post.ID, &post.Title, &post.Summary, &post.Body, &post.Date)
		if err != nil {
			return posts, fmt.Errorf("error in row.scan: %v", err)
		}
		posts.Posts = append(posts.Posts, post)
	}
	if err := rows.Close(); err != nil {
		return posts, err
	}
	return posts, nil
}

// create one post.
func (p DBPosts) createPost(post *Post, db *sql.DB) error {
	_, err := db.Exec(INSERTPOST, post.Title, post.Summary, post.Body)
	return err
}

// delete one post.
func (p DBPosts) deletePost(id string, db *sql.DB) error {
	delTime := time.Now().Format("2006-01-02 15:04:05")
	_, err := db.Exec(DELETEPOST, delTime, id)
	return err
}

// update one post.
func (p *DBPosts) updatePost(post *Post, db *sql.DB) error {
	_, err := db.Exec(UPDATEPOST, post.Title, post.Summary, post.Body, post.ID)
	return err
}
