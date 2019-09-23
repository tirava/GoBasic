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
	DSN          = "/blog?charset=utf8&interpolateParams=true"
	GETALLPOSTS  = "SELECT id, title, summary, body, updated FROM posts ORDER BY id DESC"
	GETONEPOST   = "SELECT id, title, summary, body, updated FROM posts WHERE id = ?"
)
