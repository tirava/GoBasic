package routers

import (
	"GoBasic/homework-5/myBlogBeeGo/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
}
