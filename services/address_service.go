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
	err := common.DB.Create(address).Error
	return err
}

func (s AddressService) GetAll() (*[]models.Address, error) {
	var list *[]models.Address
	err := common.DB.Model(&models.Address{}).Order("id DESC").Find(&list).Error

	return list, err
}

func (s AddressService) Update(address *models.Address) error {
	err := common.DB.Updates(address).Error
	return err
}

func (s AddressService) DeleteAll() error {
	var address *[]models.Address
	address, _ = s.GetAll()

	err := common.DB.Model(&models.Address{}).Delete(address).Error
	return err
}

func (s AddressService) DeleteById(id int) error {
	// 无效
	//err := common.DB.Where("id = ?", id).Delete(&models.Address{}).Error
	//err := common.DB.Delete(&models.Address{}, id).Error
	// 报错：WHERE conditions required
	//err := common.DB.Delete(&models.Address{Model: gorm.Model{ID: uint(id)}}).Error

	err := common.DB.Where("ID = ?", id).Delete(&models.Address{}).Error
	return err
}

func (s AddressService) Count() (int64, error) {
	var count int64
	err := common.DB.Model(&models.Address{}).Count(&count).Error

	return count, err
}

func (s AddressService) FindByProvince(province string, address *[]models.Address) error {
	return common.DB.Where("province = ?", province).Find(&address).Error
}
