package controllers

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
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

type Total struct {
	Total int `json:"total"`
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

func CSVCatatanPenjualan(c *gin.Context) {
	db, err := gorm.Open("sqlite3", dbPath)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// contactsData := c.PostForm("data")
	var total Total
	var productIn []m.ProductIn

	db.Find(&productIn)

	b := &bytes.Buffer{}
	w := csv.NewWriter(b)
	var header []string

	header = append(header, "WAKTU")
	header = append(header, "SKU")
	header = append(header, "Nama Barang")
	header = append(header, "Jumlah Pemesanan")
	header = append(header, "Jumlah Diterima")
	header = append(header, "Harga Beli")
	header = append(header, "Total")
	header = append(header, "Catatan")

	if err := w.Write(header); err != nil {
		log.Fatalln("error writing record to csv:", err)
	}

	for _, product := range productIn {
		var record []string
		total.Total = product.TotalReceived * product.PurchasePrice
		record = append(record, product.Time)
		record = append(record, product.SKU)
		record = append(record, product.Name)
		record = append(record, strconv.Itoa(product.OrderAmount))
		record = append(record, strconv.Itoa(product.TotalReceived))
		record = append(record, strconv.Itoa(product.PurchasePrice))
		record = append(record, strconv.Itoa(total.Total))
		record = append(record, product.Note)
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}
	w.Flush()

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename=catatan_barang_masuk.csv")
	c.Data(http.StatusOK, "text/csv", b.Bytes())
}

func CSVLaporanNilaiBarang(c *gin.Context) {
	db, err := gorm.Open("sqlite3", dbPath)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// contactsData := c.PostForm("data")
	t := time.Now()
	t.String()
	t.Format("2006-01-02 15:04:05")

	// var total Total
	var product []m.Product
	total := 0
	db.Exec("SELECT sum(stock) AS stocks FROM products")

	for _, product := range product {
		total = total + product.Stock
	}

	db.Find(&product)

	b := &bytes.Buffer{}
	w := csv.NewWriter(b)

	var laporanBarang []string
	laporanBarang = append(laporanBarang, "LAPORAN BARANG")
	w.Write(laporanBarang)

	var date_print []string
	var header []string
	var total_sku []string

	date_print = append(date_print, "Tanggal Cetak")
	date_print = append(date_print, fmt.Sprintf(t.Format("2006-01-02 15:04:05")))

	if err := w.Write(date_print); err != nil {
		log.Fatalln("error writing record to csv:", err)
	}

	total_sku = append(total_sku, "Total SKU")
	total_sku = append(total_sku, strconv.Itoa(len(product)))

	if err := w.Write(total_sku); err != nil {
		log.Fatalln("error writing record to csv:", err)
	}

	jumlah_total = append(total_sku, "Jumlah Total Barang")
	total_sku = append(total_sku, strconv.Itoa(len(product)))
	if err := w.Write(total_sku); err != nil {
		log.Fatalln("error writing record to csv:", err)
	}

	header = append(header, "SKU")
	header = append(header, "Nama Item")
	header = append(header, "Jumlah")
	header = append(header, "Rata-Rata Harga Beli")
	header = append(header, "Total")

	if err := w.Write(header); err != nil {
		log.Fatalln("error writing record to csv:", err)
	}

	for _, product := range product {
		var record []string
		record = append(record, product.SKU)
		record = append(record, product.Name)
		record = append(record, strconv.Itoa(product.Stock))
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}
	w.Flush()

	if err := w.Error(); err != nil {
		log.Fatal(err)
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename=Laporan_Nilai_Barang.csv")
	c.Data(http.StatusOK, "text/csv", b.Bytes())
}
