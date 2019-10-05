package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["myBlogBeeGo/controllers:APIController"] = append(beego.GlobalControllerRouter["myBlogBeeGo/controllers:APIController"],
        beego.ControllerComments{
            Method: "CreatePost",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["myBlogBeeGo/controllers:APIController"] = append(beego.GlobalControllerRouter["myBlogBeeGo/controllers:APIController"],
        beego.ControllerComments{
            Method: "GetOnePost",
            Router: `/:id([0-9a-h]+)`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["myBlogBeeGo/controllers:APIController"] = append(beego.GlobalControllerRouter["myBlogBeeGo/controllers:APIController"],
        beego.ControllerComments{
            Method: "DeletePost",
            Router: `/:id([0-9a-h]+)`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["myBlogBeeGo/controllers:APIController"] = append(beego.GlobalControllerRouter["myBlogBeeGo/controllers:APIController"],
        beego.ControllerComments{
            Method: "UpdatePost",
            Router: `/:id([0-9a-h]+)`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
