package models

import (
	"github.com/jinzhu/gorm"
)

type Tb_Barang struct {
	gorm.Model
	SKU   string `json:"sku"; gorm:"type:varchar(100); unique_index"`
	Name  string `json:"name"`
	Total int    `json:"total"`
}
