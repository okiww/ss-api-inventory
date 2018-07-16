package main

import (
	"bufio"
	"encoding/csv"
	"flag"
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
	items        m.Tb_Barang
)

func init() {
	flag.BoolVar(&runMigration, "migrate", true, "run db migration before starting the server")
	flag.BoolVar(&runSeeder, "seed", false, "run db seeder after db migration")
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
	if runMigration {
		runDBMigration()
	}

	if runSeeder {
		runDBSeeder()
	}
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}

func runDBMigration() {
	db, err := gorm.Open("sqlite3", "file:salestock.sqlite?cache=shared&mode=rwc")
	if err != nil {
		panic("failed to connect database")
	}

	// db.SetMaxOpenConns(1)
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&items)
}

func runDBSeeder() {
	db, err := gorm.Open("sqlite3", "file:salestock.sqlite?cache=shared&mode=rwc")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	csvFile, _ := os.Open("jumlah_barang.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))
	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		total, _ := strconv.Atoi(line[2])
		if line[0] != "SKU" {
			db.Create(&m.Tb_Barang{
				SKU:   line[0],
				Name:  line[1],
				Total: total,
			})
		}
	}
}
