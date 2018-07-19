package controllers

import (
	"net/http"
	m "ss-api-inventory/models"
	"strconv"
	"time"

	"github.com/biezhi/gorm-paginator/pagination"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jinzhu/gorm"
)

type newProduct struct {
	Name          string `json:"name"`
	OrderAmount   int    `json:"order_amount"`
	TotalReceived int    `json:"total_received"`
	PurchasePrice int    `json:"purchase_price"`
	TotalPrice    int    `json:"total_price"`
	ReceiptNumber string `json:"receipt_number"`
	Note          string `json:"note"`
	Time          string `json:"time"`
	SizeOfItem    string `json:"sizeOfItem"`
	Color         string `json:"color"`
}

type ProductInStruct struct {
	TotalReceived int  `json:"total_received"`
	Status        bool `json:"status"`
}

func GetProductIn(c *gin.Context) {
	db, err := gorm.Open("sqlite3", dbPath)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var products []m.ProductIn

	db.Find(&products)

	if len(products) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No products found!"})
		return
	}

	//paginate data
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	paginator := pagination.Pagging(&pagination.Param{
		DB:      db,
		Page:    page,
		Limit:   limit,
		OrderBy: []string{"id desc"},
		ShowSQL: true,
	}, &products)

	c.JSON(http.StatusOK, paginator)
}

func StoreNewProduct(c *gin.Context) {

	db, err := gorm.Open("sqlite3", dbPath)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var req newProduct
	var product m.Product
	var productIn m.ProductIn

	if c.ShouldBindWith(&req, binding.JSON) == nil {
		sku := ""

		if err := db.Where("name = ?", req.Name).First(&product).Error; err == nil {
			sku = product.SKU

			if err := db.Where("sku = ? AND status = ?", sku, 0).Find(&productIn).Error; err == nil {
				db.Model(&productIn).Updates(ProductInStruct{TotalReceived: productIn.TotalReceived + req.TotalReceived, Status: true})

				c.JSON(http.StatusCreated, gin.H{
					"status":      http.StatusCreated,
					"message":     "Product update successfully!",
					"product-sku": productIn.SKU,
				})
			} else {
				store := m.ProductIn{
					Name:          req.Name,
					OrderAmount:   req.OrderAmount,
					TotalReceived: req.TotalReceived,
					PurchasePrice: req.PurchasePrice,
					TotalPrice:    req.PurchasePrice,
					ReceiptNumber: req.ReceiptNumber,
					Time:          req.Time,
					SizeOfItem:    req.SizeOfItem,
					Color:         req.Color,
					SKU:           sku,
				}

				now := time.Now()
				store.CreatedAt = now

				db.Create(&store)

				c.JSON(http.StatusCreated, gin.H{
					"status":  http.StatusCreated,
					"message": "Product created successfully!",
				})

			}
		}

		if err := db.Where("name = ?", req.Name).First(&product).Error; err != nil {

			store := m.ProductIn{
				Name:          req.Name,
				OrderAmount:   req.OrderAmount,
				TotalReceived: req.TotalReceived,
				PurchasePrice: req.PurchasePrice,
				TotalPrice:    req.PurchasePrice,
				ReceiptNumber: req.ReceiptNumber,
				Time:          req.Time,
				SizeOfItem:    req.SizeOfItem,
				Color:         req.Color,
				SKU:           sku,
			}

			now := time.Now()
			store.CreatedAt = now

			db.Create(&store)

			c.JSON(http.StatusCreated, gin.H{
				"status":  http.StatusCreated,
				"message": "Product created successfully!",
			})

		}
	}
}
