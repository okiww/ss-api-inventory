package models

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
)

type ProductIn struct {
	gorm.Model
	//has one Product
	Product       Product `gorm:"association_foreignkey:SKU"`
	SKU           string  `json:"sku";gorm:"type:varchar(100);"`
	Name          string  `json:"name"`
	OrderAmount   int     `json:"order_amount"`
	TotalReceived int     `json:"total_received"`
	PurchasePrice int     `json:"purchase_price"`
	TotalPrice    int     `json:"total_price"`
	ReceiptNumber string  `json:"receipt_number";gorm:"type:varchar(100);`
	Status        bool    `json:"status";gorm:"default:0"`
	Note          string  `json:"note";gorm:"type:varchar(255);`
	Time          string  `json:"time"`
	SizeOfItem    string  `json:"sizeOfItem"`
	Color         string  `json:"color"`
}

func (m *ProductIn) BeforeCreate(scope *gorm.Scope) error {
	if m.SKU == "" {
		code := "SSI-D0"
		t := time.Now().Unix()                          //unix time
		a := strconv.FormatInt(t, 10)                   //format to 10 string
		sizeAcronym := AcronymSizeProduct(m.SizeOfItem) //set acronym size
		colorAcronym := AcronymColorProduct(m.Color)    // set acronym color
		randomIntgr := rand.Intn(10000)                 // set random 4 digit

		sku := code + a[:3] + strconv.Itoa(randomIntgr) + "-" + sizeAcronym + "-" + colorAcronym

		scope.SetColumn("SKU", sku)
	}
	return nil
}

func (m *ProductIn) AfterCreate(db *gorm.DB) (err error) {

	var product Product

	if err := db.Where("sku = ?", m.SKU).Find(&product).Error; err == nil {
		stock := product.Stock + m.TotalReceived
		db.Model(&product).Update("stock", stock)
	}
	return nil
}
