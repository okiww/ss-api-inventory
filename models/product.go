package models

import (
	"math/rand"
	"strconv"
	"strings"
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

//acronym color example color White convert to WH
func AcronymColorProduct(text string) string {
	words := strings.Split(text, " ")

	res := ""

	if len(words) > 1 {
		res = res + string(words[0][0]) + string(words[1][0]) + string(words[1][1])
	} else {
		res = res + string(words[0][0]) + string(words[0][1]) + string(words[0][2])
	}

	return strings.ToUpper(res)
}

//acronym size example size S convert to SS
func AcronymSizeProduct(text string) string {
	res := ""

	switch text {
	case "S":
		res = "SS"
	case "M":
		res = "MM"
	case "L":
		res = "LL"
	case "XL":
		res = "XL"
	case "XXL":
		res = "XX"
	case "XXXL":
		res = "X3"
	}

	return res
}
