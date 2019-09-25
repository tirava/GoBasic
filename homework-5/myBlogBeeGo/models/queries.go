/*
 * HomeWork-5: Start BeeGo
 * Created on 26.09.19 19:11
 * Copyright (c) 2019 - Eugene Klimov
 */

package models

import "github.com/astaxie/beego"

// DBQueries templates.
type DBQueries struct {
	QGetAllPosts string
	QGetOnePost  string
	QDeletePost  string
	QInsertPost  string
	QUpdatePost  string
}

// NewDBQueries creates new queries.
func NewDBQueries() *DBQueries {
	table := beego.AppConfig.String("TABLENAME")
	if table == "" {
		table = "posts"
	}
	return &DBQueries{
		QGetAllPosts: "SELECT id, title, summary, body, DATE_FORMAT(updated, '%d.%m.%Y %H:%i') FROM " + table + " WHERE deleted IS NULL ORDER BY id DESC",
		QGetOnePost:  "SELECT id, title, summary, body, DATE_FORMAT(updated, '%d.%m.%Y %H:%i') FROM " + table + " WHERE deleted IS NULL AND id = ?",
		QDeletePost:  "UPDATE " + table + " SET deleted = ? WHERE id = ?",
		QInsertPost:  "INSERT INTO " + table + " (title, summary, body) VALUES(?, ?, ?)",
		QUpdatePost:  "UPDATE " + table + " SET title = ?, summary = ?, body = ? WHERE ID = ?",
	}
}
