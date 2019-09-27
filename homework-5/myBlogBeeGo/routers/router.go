/*
 * HomeWork-5: Start BeeGo
 * Created on 26.09.19 22:11
 * Copyright (c) 2019 - Eugene Klimov
 */

package routers

import (
	"github.com/astaxie/beego"
	"myBlog/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{}, "get:GetAllPosts")
	beego.Router("/posts", &controllers.MainController{}, "get:GetOnePost")
	beego.Router("/posts/edit", &controllers.MainController{}, "get:GetEditPost")
	beego.Router("/posts/create", &controllers.MainController{}, "get:GetCreatePost")

	beego.Router("/api/v1/posts/:id([0-9]+", &controllers.MainController{}, "delete:DeletePost")
	beego.Router("/api/v1/posts/:id([0-9]+", &controllers.MainController{}, "put:UpdatePost")
	beego.Router("/api/v1/posts/create", &controllers.MainController{}, "post:CreatePost")
}
