package routers

import (
	"GoBasic/homework-5/myBlogBeeGo/controllers"
	"github.com/astaxie/beego"
)

const (
	POSTSURL = "/posts"
	APIURL   = "/api/v1"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router(POSTSURL, &controllers.MainController{}, "get:GetPosts")
	beego.Router(APIURL+POSTSURL+"/:id([0-9]+", &controllers.MainController{}, "delete:DeletePost")
}
