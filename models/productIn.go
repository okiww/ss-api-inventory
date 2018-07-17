package models

import (
	"github.com/jinzhu/gorm"
)

type ProductIn struct {
	gorm.Model
	//has one Product
	Product       Product `gorm:"association_foreignkey:SKU"`
	SKU           string  `json:"sku";gorm:"type:varchar(100);"`
	Name          string  `json:"name"`
	OrderAmount   int     `json:"order_amount"`
	PurchasePrice int     `json:"purchase_price"`
	TotalPrice    int     `json:"total_price"`
	ReceiptNumber string  `json:"receipt_number";gorm:"type:varchar(100);`
	Status        bool    `json:"status";gorm:"default:0"`
	Note          string  `json:"note";gorm:"type:varchar(255);`
}
