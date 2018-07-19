package models

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
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

	split_time := strings.Split(m.Time, " ")
	received := strconv.Itoa(m.TotalReceived)
	note := split_time[0] + " terima" + received
	status := setStatus(m.OrderAmount, m.TotalReceived)

	// d := time.Date(2000, 2, 1, 12, 30, 0, 0, time.UTC)
	// year, month, day := d.Date()
	// receiptNumber := year+"-"+month"-"+day-

	if status == false {
		note = split_time[0] + " terima" + received + " Menunggu"
	}
	if m.SKU == "" {
		code := "SSI-D0"
		rand.Seed(time.Now().UnixNano())
		randcode := rand.Intn(9000000)

		sizeAcronym := ""
		colorAcronym := ""

		if m.SizeOfItem != "" && m.Color != "" {
			sizeAcronym = AcronymSizeProduct(m.SizeOfItem) //set acronym size
			colorAcronym = AcronymColorProduct(m.Color)    // set acronym color
		} else {
			split1 := strings.Split(m.Name, "(")
			split2 := strings.Split(split1[1], ",")
			size := split2[0]
			split_color := strings.Split(split2[1], ")")
			color := split_color[0]

			sizeAcronym = size
			colorAcronym = color
		}
		// randomIntgr := rand.Intn(10000) // set random 4 digit

		sku := code + strconv.Itoa(randcode) + "-" + sizeAcronym + "-" + colorAcronym
		scope.SetColumn("SKU", sku)
	}
	scope.SetColumn("Note", note)
	scope.SetColumn("Status", status)
	scope.SetColumn("Status", status)
	return nil
}

func (m *ProductIn) AfterCreate(db *gorm.DB) (err error) {
	var product Product
	fmt.Println(m.SKU)

	// tx := db.Begin()
	if err := db.Where("sku = ?", m.SKU).First(&product).Error; err == nil {
		fmt.Println("update")

		stock := product.Stock + m.TotalReceived
		db.Model(&product).Update("stock", stock)
	} else {
		fmt.Println("create")
		db.Create(&Product{
			SKU:   m.SKU,
			Name:  m.Name,
			Stock: m.TotalReceived,
		})
	}

	// 	// tx.Commit()
	return nil
}

func setStatus(order_amount int, total_received int) bool {
	status := true
	if order_amount > total_received {
		status = false
	}

	return status
}
