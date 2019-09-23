/*
 * HomeWork-4: Simple blog - MySQL
 * Created on 23.09.2019 19:33
 * Copyright (c) 2019 - Eugene Klimov
 */

package main

import (
	"database/sql"
	"fmt"
	"html/template"
)

// Post is the base post type.
type Post struct {
	ID      string
	Title   string        `json:"title"`
	Date    string        `json:"date"` // todo convert to DateTime
	Summary string        `json:"summary"`
	Body    template.HTML `json:"body"`
}

//type dbPosts map[string]Post
type dbPosts []Post

func (p *dbPosts) getPosts(id string, db *sql.DB) (dbPosts, error) {
	var rows *sql.Rows
	var err error
	var posts = dbPosts{}
	if id != "" {
		rows, err = db.Query(GETONEPOST, id)
	} else {
		rows, err = db.Query(GETALLPOSTS)
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
		posts = append(posts, post)
	}
	if err := rows.Close(); err != nil {
		return posts, err
	}
	return posts, nil
}

func (p dbPosts) delete(id string) error {
	//if _, ok := p[id]; !ok {
	//	return fmt.Errorf("post not found: %v", id)
	//}
	//delete(p, id)
	return nil
}

func (p dbPosts) update(post *Post) error {
	//if _, ok := p[post.ID]; !ok {
	//	return fmt.Errorf("post not found: %v", post.ID)
	//}
	//p[post.ID] = *post
	return nil
}

func (p dbPosts) create(post *Post) error {
	//p[post.ID] = *post
	return nil
}
