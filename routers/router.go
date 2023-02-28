// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"package-example/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		// 下面的必须放在 gopath 下面的 src 目录中才会有效，否则无法生成 commentsRouter_controllers.go
		//beego.NSNamespace("/user",
		//	beego.NSInclude(
		//		&controllers.UserController{},
		//	),
		//),
		beego.NSRouter("/user/index", &controllers.UserController{}, "get:GetAll"),
		beego.NSRouter("/gjson/index", &controllers.GjsonController{}, "get:Index"),
	)
	beego.AddNamespace(ns)
}
