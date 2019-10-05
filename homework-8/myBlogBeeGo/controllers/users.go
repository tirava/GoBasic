/*
 * HomeWork-8: Config, Logs and Auth
 * Created on 06.10.19 12:16
 * Copyright (c) 2019 - Eugene Klimov
 */

package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"myBlogBeeGo/models"
	"net/http"
)

// UsersController for operations with users via API.
type UsersController struct {
	beego.Controller
}

// GetUser check user is exists.
// @Title GetUser
// @Description get user
// @Tags users
// @Param	id	path string	true	"name of user"
// @Param	body	body models.User	true	"json user body"
// @Success 200 {object} models.User
// @Failure 500 server error
// @Failure 404 not found
// @router /:id([0-9a-zA-Z]+) [post]
func (c *UsersController) GetUser() {
	userID := c.Ctx.Input.Param(":id")
	users := models.NewUser()
	user, err := c.decodeUser()
	if err != nil {
		users.Lg.Error("error while decoding new user body: %s", err)
		users.SendError(c.Ctx.ResponseWriter, http.StatusInternalServerError, err, "sorry, error while decoding new user body")
		return
	}
	if userID != user.Name {
		users.Lg.Error("userID: %s not equal userName: %s", userID, user.Name)
		users.SendError(c.Ctx.ResponseWriter, http.StatusNotFound, nil, "sorry, userID not equal userName")
		return
	}
	users.User = *user
	if err = users.GetUser(); err != nil {
		users.Lg.Error("error find user: %s", err)
		users.SendError(c.Ctx.ResponseWriter, http.StatusNotFound, err, "sorry, error find user")
		return
	}
	c.Ctx.ResponseWriter.WriteHeader(http.StatusOK)
	//c.isAuth = true
}

// CreateUser creates new user.
// @Title CreateUser
// @Description create new user
// @Tags users
// @Param	body	body models.User	true	"json user body"
// @Success 201 body is empty
// @Failure 500 server error
// @router / [post]
func (c *UsersController) CreateUser() {
	users := models.NewUser()
	user, err := c.decodeUser()
	if err != nil {
		users.Lg.Error("error while decoding new user body: %s", err)
		users.SendError(c.Ctx.ResponseWriter, http.StatusInternalServerError, err, "sorry, error while decoding new user body")
		return
	}
	users.User = *user
	if err = users.CreateUser(); err != nil {
		users.Lg.Error("error create new user: %s", err)
		users.SendError(c.Ctx.ResponseWriter, http.StatusInternalServerError, err, "sorry, error while create new user")
		return
	}
	c.Ctx.ResponseWriter.WriteHeader(http.StatusCreated)
}

// decodePost is JSON decoder helper
func (c *UsersController) decodeUser() (*models.User, error) {
	user := &models.User{}
	if err := json.NewDecoder(c.Ctx.Request.Body).Decode(user); err != nil {
		return nil, err
	}
	return user, nil
}
