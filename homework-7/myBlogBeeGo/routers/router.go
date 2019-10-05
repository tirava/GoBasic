/*
 * HomeWork-7: Testing & Docs in BeeGo
 * Created on 28.09.19 22:21
 * Copyright (c) 2019 - Eugene Klimov
 */
// @APIVersion 1.0.0
// @Title myBlog swagger API
// @Description My Blog has a cool swagger for API
// @Contact kirk@gmail.com
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html

package routers

import (
	"myBlogBeeGo/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.FormsController{}, "get:GetAllPosts")
	beego.Router("/posts", &controllers.FormsController{}, "get:GetOnePost")
	beego.Router("/posts/edit", &controllers.FormsController{}, "get:GetEditPost")
	beego.Router("/posts/create", &controllers.FormsController{}, "get:GetCreatePost")

	//beego.Router("/api/v1/posts/:id([0-9a-h]+", &controllers.ApiController{}, "delete:DeletePost")
	//beego.Router("/api/v1/posts/:id([0-9a-h]+", &controllers.ApiController{}, "put:UpdatePost")
	//beego.Router("/api/v1/posts", &controllers.ApiController{}, "post:CreatePost")

	ns := beego.NewNamespace("/api/v1/posts",
		beego.NSInclude(&controllers.APIController{}),
	)

	beego.AddNamespace(ns)
}
