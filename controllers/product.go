package controllers

import (
	"net/http"
	m "ss-api-inventory/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/jinzhu/gorm"
)

var (
	dbPath = "file:salestock.db?cache=shared&mode=rwc"
)

type (
	product struct {
		SKU       string `json:"sku"`
		Name      string `json:"name"`
		Stock     int    `json:"stock"`
		CreatedAt time.Time
	}
)

func GetProduct(c *gin.Context) {

	db, err := gorm.Open("sqlite3", dbPath)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var products []m.Product

	db.Find(&products)

	if len(products) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No users found!"})
		return
	}

	// for i, _ := range products {
	// 	db.Model(users[i]).Related(&users[i].Role)
	// }

	//transforms the todos for building a good response
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": products})
}

func CreateProduct(c *gin.Context) {
	db, err := gorm.Open("sqlite3", dbPath)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var req product
	var product m.Product
	if c.ShouldBindWith(&req, binding.JSON) == nil {

		if err := db.Where("sku = ?", req.SKU).First(&product).Error; err == nil {
			c.JSON(http.StatusOK, gin.H{"status": http.StatusNotFound, "message": "Product already in database!"})
			return
		}

		store := m.Product{
			SKU:   product.SKU,
			Name:  product.Name,
			Stock: product.Stock,
		}

		now := time.Now()
		store.CreatedAt = now

		db.Save(&store)

		c.JSON(http.StatusCreated, gin.H{
			"status":  http.StatusCreated,
			"message": "User created successfully!",
			"data":    store.Name,
		})

	}
}
