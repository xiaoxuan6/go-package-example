package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
)

type ErrorController struct {
	beego.Controller
}

func (e *ErrorController) Error404()  {
	m := map[string]interface{}{}
	m["code"] = 404
	m["msg"] = fmt.Sprintf("Method %s not found", e.Ctx.Input.URI())

	e.Data["json"] = &m
	e.ServeJSON()
}