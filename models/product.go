package models

import (
	"github.com/jinzhu/gorm"
)

type Product struct {
	gorm.Model
	SKU   string `json:"sku"; gorm:"type:varchar(100); unique_index"`
	Name  string `json:"name"`
	Stock int    `json:"stock"`
	//has many in productOut
	ProductOut []ProductOut
}
