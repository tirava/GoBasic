/*
 * HomeWork-5: Start BeeGo
 * Created on 25.09.19 23:05
 * Copyright (c) 2019 - Eugene Klimov
 */

package controllers

import (
	"GoBasic/homework-5/myBlogBeeGo/models"
	"github.com/astaxie/beego"
	"gopkg.in/russross/blackfriday.v2"
	"html/template"
	"log"
	"net/http"
)

type MainController struct {
	beego.Controller
}

const (
	BLOGNAME = "Блог Евгения Климова"
)

func (c *MainController) Get() {

	posts := models.DBPosts{DB: models.DB}
	posts, err := posts.GetPosts("")
	//log.Println(posts.Posts)
	if err != nil {
		log.Println(err)
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.Data["BlogName"] = BLOGNAME
	c.Data["Posts"] = &posts.Posts
	c.TplName = "index.tpl"
}

func (c *MainController) GetPosts() {
	postNum := c.Ctx.Request.URL.Query().Get("id")
	if postNum == "" {
		c.Redirect("/", http.StatusMovedPermanently)
	}
	posts := models.DBPosts{DB: models.DB}
	posts, err := posts.GetPosts(postNum)
	if err != nil || len(posts.Posts) == 0 {
		log.Println(err)
		c.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		c.Abort(http.StatusText(http.StatusNotFound))
	}
	c.Data["BlogName"] = BLOGNAME
	posts.Posts[0].Body = template.HTML(blackfriday.Run([]byte(posts.Posts[0].Body)))
	c.Data["Post"] = &posts.Posts[0]
	c.TplName = "post.tpl"
}
