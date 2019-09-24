/*
 * HomeWork-5: Start BeeGo
 * Created on 25.09.19 23:05
 * Copyright (c) 2019 - Eugene Klimov
 */

package controllers

import (
	"GoBasic/homework-5/myBlogBeeGo/models"
	"github.com/astaxie/beego"
	"log"
	"net/http"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {

	posts := models.DBPosts{DB: models.DB}
	posts, err := posts.GetPosts("")
	//log.Println(posts.Posts)
	if err != nil {
		log.Println(err)
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.Data["BlogName"] = "Блог Евгения Климова"
	c.Data["Posts"] = &posts.Posts
	c.TplName = "index.tpl"
}
