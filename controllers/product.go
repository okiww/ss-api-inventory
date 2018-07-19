package controllers

import (
	"net/http"
	m "ss-api-inventory/models"
	"strconv"
	"time"

	"github.com/biezhi/gorm-paginator/pagination"
	"github.com/gin-gonic/gin"
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
