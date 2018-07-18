package models

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
)

type Product struct {
	gorm.Model
	SKU        string `json:"sku"; gorm:"type:varchar(100); unique_index"`
	Name       string `json:"name"`
	Stock      int    `json:"stock"`
	SizeOfItem string `json:"sizeOfItem"`
	Color      string `json:"color"`
	//has many to productOut
	ProductOut []ProductOut
}

func (m *Product) BeforeCreate(scope *gorm.Scope) error {
	code := "SSI-D0"
	t := time.Now().Unix()                          //unix time
	a := strconv.FormatInt(t, 10)                   //format to 10 string
	sizeAcronym := AcronymSizeProduct(m.SizeOfItem) //set acronym size
	colorAcronym := AcronymColorProduct(m.Color)    // set acronym color
	randomIntgr := rand.Intn(10000)                 // set random 4 digit

	sku := code + a[:3] + strconv.Itoa(randomIntgr) + "-" + sizeAcronym + "-" + colorAcronym

	scope.SetColumn("SKU", sku)

	return nil
}
