# SS API INVENTORY

this is sample API inventory

## Getting Started

This project using gin-gonic so you must install gin-gonic first

## Installation

To install Gin package, you need to install Go and set your Go workspace first.
[gin-gonic](https://github.com/gin-gonic)

1. Download and install it:

```sh
$ go get -u github.com/gin-gonic/gin
```

2. Import it in your code:

```go
import "github.com/gin-gonic/gin"
```

3. (Optional) Import `net/http`. This is required for example if using constants such as `http.StatusOK`.

```go
import "net/http"
```

## Starting

1. Go to
```sh
$ cd $GOPATH/src/
```
2. Clone repository & run - will be run at :8080
```sh
$ go run main.go

[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /ping                     --> main.setupRouter.func1 (3 handlers)
[GIN-debug] GET    /catatan-barang-masuk     --> ss-api-inventory/controllers.CSVCatatanPenjualan (3 handlers)
[GIN-debug] GET    /laporan-nilai-barang     --> ss-api-inventory/controllers.CSVLaporanNilaiBarang (3 handlers)
[GIN-debug] GET    /api/v1/products          --> ss-api-inventory/controllers.GetProduct (3 handlers)
[GIN-debug] GET    /api/v1/product/:sku      --> ss-api-inventory/controllers.GetProductBySku (3 handlers)
[GIN-debug] DELETE /api/v1/product/:sku      --> ss-api-inventory/controllers.DeleteProduct (3 handlers)
[GIN-debug] GET    /api/v1/products/in       --> ss-api-inventory/controllers.GetProductIn (3 handlers)
[GIN-debug] POST   /api/v1/products/in       --> ss-api-inventory/controllers.StoreNewProduct (3 handlers)
[GIN-debug] Listening and serving HTTP on :8080
```

3. Set to TRUE/FALSE for run migration & seeder it will be automatically when run the program
```go
func init() {
	//for turn on migration & seeder change to true
	flag.BoolVar(&runMigration, "migrate", true, "run db migration before starting the server")
	flag.BoolVar(&runSeeder, "seed", false, "run db seeder after db migration")
	flag.Parse()
}
```

4. Generate report
Open browser and go to localhost:8080
```sh
http://localhost:8080/catatan-barang-masuk
http://localhost:8080/laporan-nilai-barang               
```

## API

List of API

API                           |   Method   |                        Note                      |
--------------------------------------------|-----------:|--------------------------------------------------|
http://localhost:8080/ping                  |    GET     | Test connection API                              |
http://localhost:8080/api/v1/products       |    GET     | Get list data Product                            |
http://localhost:8080/api/v1/product/:sku   |    GET     | Get lisr data Product by SKU                     |
http://localhost:8080/api/v1/product/:sku   |    DELETE  | Delete data Product by SKU                       |
http://localhost:8080//api/v1/products/in   |    GET     | GET list data "barang masuk"                     |
http://localhost:8080//api/v1/products/in   |    POST    | Store new items                                  |


Add additional notes about how to deploy this on a live system
## API Examples
POST DATA http://localhost:8080//api/v1/products/in
give body example like this :

```go
{
	"name": "koko",
	"order_amount": 50,
	"total_received": 50,
	"purchase_price": 10000,
	"ReceiptNumber": "123213",
	"Time": "2018/01/02 11:20",
	"sizeOfItem": "M",
	"note": "hahahaha",
	"color": "Red Green"
}
```

data automate store on salestock.db

## Built With

* [gin-gonic](https://github.com/gin-gonic) - The web framework used
* [GORM](http://gorm.io/) - The ORM used
