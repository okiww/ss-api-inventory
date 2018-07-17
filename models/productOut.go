package models

import (
	"github.com/jinzhu/gorm"
)

type ProductOut struct {
	gorm.Model
	//has one Product
	Product      Product `gorm:"association_foreignkey:SKU"`
	SKU          string  `json:"sku"; gorm:"type:varchar(100);"`
	Name         string  `json:"name"`
	NumberOfItem int     `json:"number_of_item"`
	SellingPrice int     `json:"selling_price"`
	TotalPrice   int     `json:"total_price"`
	Note         string  `json:"note"; gorm:"type:varchar(255);`
	Time         string
}
