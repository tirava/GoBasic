/*
 * HomeWork-7: Testing & Docs in BeeGo
 * Created on 28.09.19 22:17
 * Copyright (c) 2019 - Eugene Klimov
 */

package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"myBlogBeeGo/models"
	"net/http"
)

// APIController for operations with posts via API.
type APIController struct {
	beego.Controller
}

//GetOnePost shows one posts with full content.
//@Title GetOnePost
//@Description get one post
//@Tags posts
//@Param	id	path string	true	"ID of the post"
//@Success 200 {object} models.Post
//@Failure 500 body is empty
//@Failure 404 not found
//@router /:id([0-9a-h]+) [get]
func (c *APIController) GetOnePost() {
	postNum := c.Ctx.Input.Param(":id")
	posts := models.NewPosts()
	if err := posts.GetPosts(postNum); err != nil {
		posts.Lg.Error("error get one post: %s", err)
		posts.SendError(c.Ctx.ResponseWriter, http.StatusNotFound, err, "sorry, error searching post")
		return
	}
	c.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(posts.Posts[0])
	if err != nil {
		posts.Lg.Error("Can't marshal error data: %s", err)
		posts.SendError(c.Ctx.ResponseWriter, http.StatusInternalServerError, err, "sorry, error getting post")
		return
	}
	if _, err = c.Ctx.ResponseWriter.Write(data); err != nil {
		posts.Lg.Error("Can't write to ResponseWriter: %s", err)
		return
	}
}

// DeletePost removes post from DB.
// @Title DeletePost
// @Description delete post
// @Tags posts
// @Param	id	path string	true	"ID of the post"
// @Success 200 body is empty
// @Failure 500 body is empty
// @Failure 404 not found
// @router /:id([0-9a-h]+) [delete]
func (c *APIController) DeletePost() {
	postNum := c.Ctx.Input.Param(":id")
	posts := models.NewPosts()
	if err := posts.DeletePost(postNum); err != nil {
		posts.Lg.Error("error delete post: %s", err)
		posts.SendError(c.Ctx.ResponseWriter, http.StatusInternalServerError, err, "sorry, error deleting post")
		return
	}
	c.Ctx.ResponseWriter.WriteHeader(http.StatusOK)
}

// UpdatePost updates post in DB.
// @Title UpdatePost
// @Description update post
// @Tags posts
// @Param	id	path string	true	"ID of the post"
// @Param	body	body models.Post	true	"json post body, use html body in double quotes instead {}"
// @Success 200 body is empty
// @Failure 500 body is empty
// @Failure 404 not found
// @router /:id([0-9a-h]+) [put]
func (c *APIController) UpdatePost() {
	postNum := c.Ctx.Input.Param(":id")
	posts := models.NewPosts()
	post, err := c.decodePost()
	if err != nil {
		posts.Lg.Error("error while decoding post body: %s", err)
		posts.SendError(c.Ctx.ResponseWriter, http.StatusInternalServerError, err, "sorry, error while decoding post body")
		return
	}
	posts.Posts = append(posts.Posts, *post)
	if err = posts.UpdatePost(postNum, false); err != nil {
		posts.Lg.Error("error edit post: %s", err)
		posts.SendError(c.Ctx.ResponseWriter, http.StatusInternalServerError, err, "sorry, error while edit post")
		return
	}
	c.Ctx.ResponseWriter.WriteHeader(http.StatusOK)
}

// CreatePost creates new post.
// @Title CreatePost
// @Description create new post
// @Tags posts
// @Param	body	body models.Post	true	"json post body, use html body in double quotes instead {}"
// @Success 201 body is empty
// @Failure 500 body is empty
// @router / [post]
func (c *APIController) CreatePost() {
	posts := models.NewPosts()

	post, err := c.decodePost()
	if err != nil {
		posts.Lg.Error("error while decoding new post body: %s", err)
		posts.SendError(c.Ctx.ResponseWriter, http.StatusInternalServerError, err, "sorry, error while decoding new post body")
		return
	}
	posts.Posts = append(posts.Posts, *post)
	if err = posts.CreatePost(); err != nil {
		posts.Lg.Error("error create new post: %s", err)
		posts.SendError(c.Ctx.ResponseWriter, http.StatusInternalServerError, err, "sorry, error while create new post")
		return
	}
	c.Ctx.ResponseWriter.WriteHeader(http.StatusCreated)
}

// decodePost is JSON decoder helper
func (c *APIController) decodePost() (*models.Post, error) {
	post := &models.Post{}
	if err := json.NewDecoder(c.Ctx.Request.Body).Decode(post); err != nil {
		return nil, err
	}
	if post.ID != "" {
		post.OID, _ = primitive.ObjectIDFromHex(post.ID)
	}
	return post, nil
}
