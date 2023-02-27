package controllers

import (
	"package-example/models"
	"package-example/services"
)

type AddressController struct {
	BaseController
}

func (a *AddressController) Index() {

	list, err := services.NewAddressService().Index()
	if err != nil {
		a.Error(err.Error())
	}

	a.Output(list)
}

func (a *AddressController) Create() {
	address := &models.Address{
		Province: "山东省",
		City:     "临沂市",
		Area:     "兰山区",
		Address:  "金佰瀚110号",
	}

	err := services.NewAddressService().Add(address)
	if err != nil {
		a.Error(err.Error())
	}

	a.OutputMsg("添加成功")
}
