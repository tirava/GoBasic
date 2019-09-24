package routers

import (
	"GoBasic/homework-5/myBlogBeeGo/controllers"
	"github.com/astaxie/beego"
)

const (
	POSTSURL = "/posts"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router(POSTSURL, &controllers.MainController{}, "get:GetPosts")
}
