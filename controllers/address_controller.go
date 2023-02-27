package controllers

import (
	"fmt"
	"gorm.io/gorm"
	"package-example/models"
	"package-example/services"
)

type AddressController struct {
	BaseController
}

func (a *AddressController) Index() {

	list, err := services.NewAddressService().GetAll()
	if err != nil {
		a.Error(fmt.Sprintf("exec getAll err：%s", err.Error()))
	}

	count, err := services.NewAddressService().Count()
	if err != nil {
		a.Error(fmt.Sprintf("exec count err：%s", err.Error()))
	}

	item := map[string]interface{}{}
	item["list"] = list
	item["count"] = count
	a.Output(item)
}

func (a *AddressController) Create() {
	address := &models.Address{
		Province: "山东省",
		City:     "临沂市",
		Area:     "兰山区",
		Address:  "金佰瀚110号",
	}

	if err := services.NewAddressService().Add(address); err != nil {
		a.Error(err.Error())
	}

	a.OutputMsg("添加成功")
}

func (a *AddressController) DeleteAll() {

	if err := services.NewAddressService().DeleteAll(); err != nil {
		a.Error(err.Error())
	}

	a.OutputMsg("删除成功")
}

func (a AddressController) DeleteById() {
	id, _ := a.GetInt("id")
	if err := services.NewAddressService().DeleteById(id); err != nil {
		a.Error(fmt.Sprintf("exec delete by id err: %s", err.Error()))
	}

	a.OutputMsg("删除成功")
}

func (a *AddressController) Update() {
	address := &models.Address{
		Model: gorm.Model{
			ID: 1,
		},
		Province: "陕西省",
	}

	if err := services.NewAddressService().Update(address); err != nil {
		a.Error(err.Error())
	}

	a.OutputMsg("修改成功")

}

func (a AddressController) Search() {
	address := &[]models.Address{}
	if err := services.NewAddressService().FindByProvince("山东省", address); err != nil {
		a.Error(fmt.Sprintf("exec find err：%s", err.Error()))
	}

	a.Output(address)

}
