/*
 * HomeWork-4: Simple blog - MySQL
 * Created on 23.09.2019 21:33
 * Copyright (c) 2019 - Eugene Klimov
 */

package main

// Constants
const (
	SERVADDR     = ":8080"
	TEMPLATEEXT  = "*.gohtml"
	TEMPLATEPATH = "templates"
	POSTSURL     = "/posts"
	EDITURL      = "/edit"
	CREATEURL    = "/create"
	APIURL       = "/api/v1"
	STATICPATH   = "/static"
	DBNAME       = "blog"
	DSN          = "/" + DBNAME + "?charset=utf8&interpolateParams=true"
	TABLENAME    = "posts"
	GETALLPOSTS  = "SELECT id, title, summary, body, updated FROM " + TABLENAME + " ORDER BY id DESC"
	GETONEPOST   = "SELECT id, title, summary, body, updated FROM " + TABLENAME + " WHERE id = ?"
	DELETEPOST   = "DELETE FROM " + TABLENAME + " WHERE id = ?"
	INSERTPOST   = "INSERT INTO " + TABLENAME + " (title, summary, summary) VALUES(?, ?, ?)"
	UPDATEPOST   = "UPDATE " + TABLENAME + " SET title = ?, summary = ? summary = ? WHERE ID = ?"
)
