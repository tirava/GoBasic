/*
 * HomeWork-7: Testing & Docs in BeeGo
 * Created on 04.10.19 19:17
 * Copyright (c) 2019 - Eugene Klimov
 */

package controllers

import (
	"html/template"
	"myBlogBeeGo/models"
	"net/http"

	"github.com/astaxie/beego"
	"gopkg.in/russross/blackfriday.v2"
)

// FormsController for operations with posts via Forms (get one, edit, create).
type FormsController struct {
	beego.Controller
}

// GetAllPosts shows all posts in main page.
func (c *FormsController) GetAllPosts() {
	posts := models.NewPosts()
	if err := posts.GetPosts(""); err != nil {
		posts.Lg.Error("error get all posts: %s", err)
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Abort("500")
		return
	}
	c.Data["BlogName"] = beego.AppConfig.String("BLOGNAME")
	c.Data["Posts"] = &posts.Posts
	c.TplName = "index.tpl"
}

// GetOnePost shows one posts with full content.
func (c *FormsController) GetOnePost() {
	postNum := c.Ctx.Request.URL.Query().Get("id")
	if postNum == "" {
		c.Redirect("/", http.StatusMovedPermanently)
	}
	posts := models.NewPosts()
	if err := posts.GetPosts(postNum); err != nil {
		posts.Lg.Error("error get one post: %s", err)
		c.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		c.Abort("404")
	}
	c.Data["BlogName"] = beego.AppConfig.String("BLOGNAME")
	posts.Posts[0].Body = template.HTML(blackfriday.Run([]byte(posts.Posts[0].Body)))
	posts.Posts[0].ID = posts.Posts[0].OID.Hex()
	c.Data["Post"] = &posts.Posts[0]
	c.TplName = "post.tpl"
}

// GetEditPost shows edit form for edit post.
func (c *FormsController) GetEditPost() {
	postNum := c.Ctx.Request.URL.Query().Get("id")
	if postNum == "" {
		c.Redirect("/", http.StatusMovedPermanently)
	}
	posts := models.NewPosts()
	if err := posts.GetPosts(postNum); err != nil {
		posts.Lg.Error("error get one post for edit: %s", err)
		c.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		c.Abort("404")
	}
	c.Data["BlogName"] = beego.AppConfig.String("BLOGNAME")
	posts.Posts[0].ID = posts.Posts[0].OID.Hex()
	c.Data["Post"] = &posts.Posts[0]
	c.TplName = "edit.tpl"
}

// GetCreatePost shows clean form for new post.
func (c *FormsController) GetCreatePost() {
	c.Data["BlogName"] = beego.AppConfig.String("BLOGNAME")
	c.TplName = "create.tpl"
}
