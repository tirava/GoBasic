package routers

import (
	"GoBasic/homework-5/myBlogBeeGo/controllers"
	"github.com/astaxie/beego"
)

// Constants.
const (
	POSTSURL  = "/posts"
	APIURL    = "/api/v1"
	EDITURL   = "/edit"
	CREATEURL = "/create"
)

func init() {
	beego.Router("/", &controllers.MainController{}, "get:GetPosts")
	beego.Router(POSTSURL, &controllers.MainController{}, "get:GetPost")
	beego.Router(POSTSURL+EDITURL, &controllers.MainController{}, "get:GetEditPost")
	beego.Router(POSTSURL+CREATEURL, &controllers.MainController{}, "get:GetCreatePost")

	beego.Router(APIURL+POSTSURL+"/:id([0-9]+", &controllers.MainController{}, "delete:DeletePost")
	beego.Router(APIURL+POSTSURL+"/:id([0-9]+", &controllers.MainController{}, "put:UpdatePost")
	beego.Router(APIURL+POSTSURL+CREATEURL, &controllers.MainController{}, "post:CreatePost")
}
