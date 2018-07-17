package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	ctrl "ss-api-inventory/controllers"
	m "ss-api-inventory/models"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	runMigration bool
	runSeeder    bool
	product      m.Product
	productIn    m.ProductIn
	productOut   m.ProductOut
	dbPath       = "file:salestock.db?cache=shared&mode=rwc"
)

func init() {
	//for turn on migration & seeder change to true
	flag.BoolVar(&runMigration, "migrate", true, "run db migration before starting the server")
	flag.BoolVar(&runSeeder, "seed", false, "run db seeder after db migration")
	flag.Parse()
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	v1 := r.Group("/api/v1")
	// v1.Use(db)
	v1.GET("/products", ctrl.GetProduct)
	v1.GET("/product/:sku", ctrl.GetProductBySku)
	v1.POST("/product", ctrl.CreateProduct)
	v1.DELETE("/product/:sku", ctrl.DeleteProduct)

	return r
}

func main() {
	db, err := getDatabaseHandle()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if runMigration {
		runDBMigration(db)
	}

	if runSeeder {
		runDBSeeder(db)
	}
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}

func runDBMigration(db *gorm.DB) {
	// Migrate the schema
	db.AutoMigrate(&product, &productOut, &productIn)
}

//Feature Task : Optional: Import data from CSV/spreadsheet (data migration)
func runDBSeeder(db *gorm.DB) {
	tx := db.Begin()

	product_file, _ := os.Open("jumlah_barang.csv")
	product_in, _ := os.Open("catatan_barang_masuk.csv")
	product_out, _ := os.Open("catatan_barang_keluar.csv")

	//seeder for table product
	reader_product := csv.NewReader(bufio.NewReader(product_file))

	for {
		// fmt.Println(header)
		line, error := reader_product.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		total, _ := strconv.Atoi(line[2])
		//delete header
		if line[0] != "SKU" {
			tx.FirstOrCreate(&m.Product{
				SKU:   line[0],
				Name:  line[1],
				Stock: total,
			})
		}
	}

	// seeder from csv for table product_in
	reader_product_in := csv.NewReader(bufio.NewReader(product_in))
	for {
		line, error := reader_product_in.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}

		if line[1] != "SKU" {
			order_amount, _ := strconv.Atoi(line[3])
			total_received, _ := strconv.Atoi(line[4])
			purchase_price := stringToAmount(line[5])
			total_price := stringToAmount(line[6])
			var status bool

			if order_amount > total_received {
				status = false
			} else {
				status = true
			}

			tx.FirstOrCreate(&m.ProductIn{
				Time:          line[0],
				SKU:           line[1],
				Name:          line[2],
				OrderAmount:   order_amount,
				TotalReceived: total_received,
				PurchasePrice: purchase_price,
				TotalPrice:    total_price,
				ReceiptNumber: line[7],
				Status:        status,
				Note:          line[8],
			})
		}
	}

	//seeder for table productOut
	reader_product_out := csv.NewReader(bufio.NewReader(product_out))

	for {
		line, error := reader_product_out.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}

		if line[1] != "SKU" {
			number_of_item, _ := strconv.Atoi(line[3])
			total_out := stringToAmount(line[4])
			total_price := stringToAmount(line[5])

			tx.FirstOrCreate(&m.ProductOut{
				Time:         line[0],
				SKU:          line[1],
				Name:         line[2],
				NumberOfItem: number_of_item,
				SellingPrice: total_out,
				TotalPrice:   total_price,
				Note:         line[6],
			})
		}
	}

	tx.Commit()
}

func getDatabaseHandle() (*gorm.DB, error) {
	database, err := gorm.Open("sqlite3", dbPath)
	if err != nil {
		fmt.Println("Failed to create the handle")
		return nil, err
	}
	if err2 := database.DB().Ping(); err2 != nil {
		fmt.Println("Failed to keep connection alive")
		return nil, err
	}
	return database, nil
}

func stringToAmount(str string) int {
	charIndex := str[0:3]
	slice := str[len(charIndex)-1:]
	replaceComa := strings.Replace(slice, ",", "", -1)

	parseInt, _ := strconv.Atoi(replaceComa)

	return parseInt
}
