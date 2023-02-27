package controllers

import "github.com/astaxie/beego"

type BaseController struct {
	beego.Controller
}

func (b *BaseController) Output(data interface{}) {
	m := map[string]interface{}{}
	m["code"] = 200
	m["item"] = data
	m["msg"] = ""

	b.Data["json"] = &m
	b.ServeJSON()
}

func (b *BaseController) OutputMsg(message string) {
	m := map[string]interface{}{}
	m["code"] = 200
	m["item"] = ""
	m["msg"] = message

	b.Data["json"] = &m
	b.ServeJSON()
}

func (b *BaseController) Error(message string) {
	m := map[string]interface{}{}
	m["code"] = 500
	m["item"] = ""
	m["msg"] = message

	b.Data["json"] = &m
	b.ServeJSON()
}
