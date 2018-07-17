package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	m "ss-api-inventory/models"
	"strconv"

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
	flag.BoolVar(&runMigration, "migrate", true, "run db migration before starting the server")
	flag.BoolVar(&runSeeder, "seed", true, "run db seeder after db migration")
	flag.Parse()
}

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

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
	product_file, _ := os.Open("jumlah_barang.csv")
	// product_in, _ := os.Open("catatan_barang_masuk.csv")
	//seeder for table product
	reader_product := csv.NewReader(bufio.NewReader(product_file))
	for {
		line, error := reader_product.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		total, _ := strconv.Atoi(line[2])
		if line[0] != "SKU" {
			db.Create(&m.Product{
				SKU:   line[0],
				Name:  line[1],
				Stock: total,
			})
		}
	}

	// reader_product_in := csv.NewReader(bufio.NewReader(product_in))
	// for {
	// 	line, error := reader_product_in.Read()
	// 	if error == io.EOF {
	// 		break
	// 	} else if error != nil {
	// 		log.Fatal(error)
	// 	}

	// 	if line[0] != "SKU" {
	// 		db.Create(&m.Product{
	// 			SKU:   line[0],
	// 			Name:  line[1],
	// 			Stock: total,
	// 		})
	// 	}
	// }

	//seeder for table productOut
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
