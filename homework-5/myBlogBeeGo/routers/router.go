package routers

import (
	"GoBasic/homework-5/myBlogBeeGo/controllers"
	"github.com/astaxie/beego"
)

const (
	POSTSURL = "/posts"
	APIURL   = "/api/v1"
	EDITURL  = "/edit"
)

func init() {
	beego.Router("/", &controllers.MainController{}, "get:GetPosts")
	beego.Router(POSTSURL, &controllers.MainController{}, "get:GetPost")
	beego.Router(POSTSURL+EDITURL, &controllers.MainController{}, "get:EditPost")
	beego.Router(APIURL+POSTSURL+"/:id([0-9]+", &controllers.MainController{}, "delete:DeletePost")
	beego.Router(APIURL+POSTSURL+"/:id([0-9]+", &controllers.MainController{}, "put:UpdatePost")
}
