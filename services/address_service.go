package services

import (
	"package-example/common"
	"package-example/models"
)

type AddressService struct {
}

func NewAddressService() *AddressService {
	return new(AddressService)
}

func (s AddressService) Add(address *models.Address) error {
	err := common.DB.Model(&models.Address{}).Create(address).Error
	return err
}

func (s AddressService) Index() (*[]models.Address, error) {
	var list *[]models.Address
	err := common.DB.Model(&models.Address{}).Order("id DESC").Find(&list).Error

	return list, err
}
