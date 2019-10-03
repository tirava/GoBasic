package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["bapi/controllers:MainController"] = append(beego.GlobalControllerRouter["bapi/controllers:MainController"],
        beego.ControllerComments{
            Method: "CreatePost",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["bapi/controllers:MainController"] = append(beego.GlobalControllerRouter["bapi/controllers:MainController"],
        beego.ControllerComments{
            Method: "DeletePost",
            Router: `/:id([0-9a-h]+)`,
            AllowHTTPMethods: []string{"delete"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
