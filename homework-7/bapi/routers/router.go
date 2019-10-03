// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"bapi/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{}, "get:GetAllPosts")
	beego.Router("/posts", &controllers.MainController{}, "get:GetOnePost")
	beego.Router("/posts/edit", &controllers.MainController{}, "get:GetEditPost")
	beego.Router("/posts/create", &controllers.MainController{}, "get:GetCreatePost")

	beego.Router("/api/v1/posts/:id([0-9a-h]+", &controllers.MainController{}, "delete:DeletePost")
	beego.Router("/api/v1/posts/:id([0-9a-h]+", &controllers.MainController{}, "put:UpdatePost")
	//beego.Router("/api/v1/posts", &controllers.MainController{}, "post:CreatePost")

	ns := beego.NewNamespace("api/v1/posts",
		beego.NSRouter("", &controllers.MainController{}, "post:CreatePost"))

	beego.AddNamespace(ns)
}

/*
 * HomeWork-7: Testing & Docs in BeeGo
 * Created on 28.09.19 22:21
 * Copyright (c) 2019 - Eugene Klimov
 */
