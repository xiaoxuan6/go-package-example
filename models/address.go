package models

import "gorm.io/gorm"

type Address struct {
	gorm.Model
	Province string `type:"varchar(32); not null;comment:'省份'" json:"province"`
	City     string `type:"varchar(32); not null;comment:'市'" json:"city"`
	Area     string `type:"varchar(64); not null;comment:'区'" json:"area"`
	Address  string `type:"varchar(255); not null;comment:'详细地址'" json:"address"`
}
