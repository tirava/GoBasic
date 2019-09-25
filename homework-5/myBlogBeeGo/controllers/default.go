/*
 * HomeWork-5: Start BeeGo
 * Created on 25.09.19 23:05
 * Copyright (c) 2019 - Eugene Klimov
 */

package controllers

import (
	"GoBasic/homework-5/myBlogBeeGo/models"
	"encoding/json"
	"github.com/astaxie/beego"
	"gopkg.in/russross/blackfriday.v2"
	"html/template"
	"net/http"
)

type MainController struct {
	beego.Controller
}

const (
	BLOGNAME = "Блог Евгения Климова"
)

// GetPosts shows all posts in main page.
func (c *MainController) GetPosts() {
	posts := models.NewPosts()
	err := posts.GetPosts("")
	if err != nil {
		posts.Lg.Error("error get all posts: %s", err)
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Abort("500")
		return
	}
	c.Data["BlogName"] = BLOGNAME
	c.Data["Posts"] = &posts.Posts
	c.TplName = "index.tpl"
}

// GetPost shows one posts with full content.
func (c *MainController) GetPost() {
	postNum := c.Ctx.Request.URL.Query().Get("id")
	if postNum == "" {
		c.Redirect("/", http.StatusMovedPermanently)
	}
	posts := models.NewPosts()
	err := posts.GetPosts(postNum)
	if err != nil {
		posts.Lg.Error("error get one post: %s", err)
		c.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		c.Abort("404")
	}
	c.Data["BlogName"] = BLOGNAME
	posts.Posts[0].Body = template.HTML(blackfriday.Run([]byte(posts.Posts[0].Body)))
	c.Data["Post"] = &posts.Posts[0]
	c.TplName = "post.tpl"
}

// DeletePost removes post from DB.
func (c *MainController) DeletePost() {
	postNum := c.Ctx.Input.Param(":id")
	posts := models.NewPosts()
	err := posts.DeletePost(postNum)
	if err != nil {
		posts.Lg.Error("error delete post: %s", err)
		posts.SendError(c.Ctx.ResponseWriter, http.StatusInternalServerError, err, "sorry, error while delete post")
		return
	}
	c.Ctx.ResponseWriter.WriteHeader(http.StatusOK)
}

// EditPost shows edit form for edit post.
func (c *MainController) EditPost() {
	postNum := c.Ctx.Request.URL.Query().Get("id")
	if postNum == "" {
		c.Redirect("/", http.StatusMovedPermanently)
	}
	posts := models.NewPosts()
	err := posts.GetPosts(postNum)
	if err != nil {
		posts.Lg.Error("error get one post for edit: %s", err)
		c.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		c.Abort("404")
	}
	c.Data["BlogName"] = BLOGNAME
	c.Data["Post"] = &posts.Posts[0]
	c.TplName = "edit.tpl"
}

// UpdatePost updates post in DB.
func (c *MainController) UpdatePost() {
	postNum := c.Ctx.Input.Param(":id")
	posts := models.NewPosts()

	post := models.Post{}
	err := json.NewDecoder(c.Ctx.Request.Body).Decode(&post)
	if err != nil {
		posts.Lg.Error("error while decoding post body: %s", err)
		posts.SendError(c.Ctx.ResponseWriter, http.StatusInternalServerError, err, "sorry, error while decoding post body")
		return
	}
	post.ID = postNum
	posts.Posts = append(posts.Posts, post)
	err = posts.UpdatePost()
	if err != nil {
		posts.Lg.Error("error edit post: %s", err)
		posts.SendError(c.Ctx.ResponseWriter, http.StatusInternalServerError, err, "sorry, error while edit post")
		return
	}
	c.Ctx.ResponseWriter.WriteHeader(http.StatusOK)
}
