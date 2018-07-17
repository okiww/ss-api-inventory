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

type product struct {
	Name       string `json:"name"`
	Stock      int    `json:"stock"`
	SizeOfItem string `json:"sizeOfItem"`
	Color      string `json:"color"`
	CreatedAt  time.Time
}

//Get All data Product
func GetProduct(c *gin.Context) {
	db, err := gorm.Open("sqlite3", dbPath)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var products []m.Product

	db.Find(&products)

	if len(products) <= 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No products found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": products})
}

//GET DATA PRODUCT BY ID
func GetProductBySku(c *gin.Context) {
	db, err := gorm.Open("sqlite3", dbPath)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var product m.Product
	sku := c.Param("sku")

	if err := db.Where("sku LIKE ?", sku).Find(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No product found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "data": product})
}

//INSERT DATA PRODUCT
func CreateProduct(c *gin.Context) {
	db, err := gorm.Open("sqlite3", dbPath)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var req product
	var product m.Product
	if c.ShouldBindWith(&req, binding.JSON) == nil {

		if err := db.Where("name = ?", req.Name).First(&product).Error; err == nil {
			c.JSON(http.StatusOK, gin.H{"status": http.StatusNotFound, "message": "Product already in database!"})
			return
		}

		store := m.Product{
			SizeOfItem: req.SizeOfItem,
			Name:       req.Name,
			Stock:      req.Stock,
			Color:      req.Color,
		}

		now := time.Now()
		store.CreatedAt = now

		db.Save(&store)

		c.JSON(http.StatusCreated, gin.H{
			"status":      http.StatusCreated,
			"message":     "Product created successfully!",
			"product-sku": store.SKU,
		})

	}
}

func DeleteProduct(c *gin.Context) {

	db, err := gorm.Open("sqlite3", dbPath)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	var product m.Product
	sku := c.Param("sku")

	if err := db.Where("sku LIKE ?", sku).Find(&product).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "No product found!"})
		return
	}

	db.Delete(&product)
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "Product deleted successfully!"})
}
