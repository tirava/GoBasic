/*
 * HomeWork-4: Simple blog - MySQL
 * Created on 23.09.2019 19:33
 * Copyright (c) 2019 - Eugene Klimov
 */

package main

import (
	"fmt"
	"html/template"
	"strconv"
)

// Post is the base post type.
type Post struct {
	ID      int
	Title   string        `json:"title"`
	Date    string        `json:"date"` // todo convert to DateTime
	Summary string        `json:"summary"`
	Body    template.HTML `json:"body"`
}

type dbPosts map[string]Post

func (p dbPosts) delete(id string) error {
	if _, ok := p[id]; !ok {
		return fmt.Errorf("post not found: %v", id)
	}
	delete(p, id)
	return nil
}

func (p dbPosts) update(post *Post) error {
	ids := strconv.Itoa(post.ID)
	if _, ok := p[ids]; !ok {
		return fmt.Errorf("post not found: %v", post.ID)
	}
	p[ids] = *post
	return nil
}

func (p dbPosts) create(post *Post) error {
	ids := strconv.Itoa(post.ID)
	p[ids] = *post
	return nil
}
