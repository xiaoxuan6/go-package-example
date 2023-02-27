package controllers

import "package-example/services"

type CacheController struct {
	BaseController
}

func (c CacheController) Get() {
	val := services.CacheService.Get("name")
	c.Output(val)
}

func (c CacheController) Set() {
	services.CacheService.Set("name", "vinhson")
	c.OutputMsg("添加成功")
}

func (c CacheController) Del() {
	services.CacheService.Del("name")
	c.OutputMsg("删除成功")
}
