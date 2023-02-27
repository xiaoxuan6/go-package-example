package main

import (
	"package-example/controllers"
	_ "package-example/routers"

	"github.com/astaxie/beego"
)

func main() {
	// 生成文档：bee run -gendoc=true -downdoc=true
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.ErrorController(&controllers.ErrorController{})
	beego.Run()
}
