/*
 * HomeWork-6: Mongo in BeeGo
 * Created on 28.09.19 22:19
 * Copyright (c) 2019 - Eugene Klimov
 */

package conf

import (
	"fmt"
	"github.com/astaxie/beego"
)

// GetDSN returns DSN from config or defaults
func GetDSN() string {
	dbName := beego.AppConfig.String("DBNAME")
	if dbName == "" {
		dbName = "blog"
	}
	dsn := beego.AppConfig.String("DSN")
	if dsn == "" {
		dsn = "?charset=utf8&interpolateParams=true"
	}
	host := beego.AppConfig.String("DBHOST")
	if host == "" {
		host = "localhost"
	}
	port := beego.AppConfig.String("DBPORT")
	if port == "" {
		port = "3306"
	}
	user := beego.AppConfig.String("DBUSER")
	if user == "" {
		user = "testuser"
	}
	pass := beego.AppConfig.String("DBPASS")

	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s", user, pass, host, port, dbName, dsn)
}
