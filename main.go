package main

import (
	"github.com/astaxie/beego"
	"package-example/common"
	"package-example/controllers"
	_ "package-example/routers"
)

func main() {
	// 生成文档：bee run -gendoc=true -downdoc=true
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	// 初始化 cache
	common.Init()
	beego.ErrorController(&controllers.ErrorController{})
	beego.Run()
}
