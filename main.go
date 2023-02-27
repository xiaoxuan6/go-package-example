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
	beego.ErrorController(&controllers.ErrorController{})
	// 初始化数据库
	common.Init()
	beego.Run()
}
