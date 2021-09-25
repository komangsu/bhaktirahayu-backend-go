package models

import (
	"restapi/config"
	"restapi/schemas"
)

// get stock
func GetStockLogistics(name string) int {
	var stock int

	db := config.InitDB()
	defer db.Close()

	query := `select stock from logistics where name = $1`
	db.QueryRow(query, name).Scan(&stock)

	return stock
}

// add new stock
func CreateNewStock(add_stock schemas.AddStockSchema) {
	db := config.InitDB()
	defer db.Close()

	query := `update logistics set stock = $1, price = $2 where name = $3`

	stmt, err := db.Prepare(query)
	if err != nil {
		panic(err)
	}
	stmt.Exec(add_stock.Stock, add_stock.Price, add_stock.Name)

	return
}

// check drug name in database
func CheckDrugName(name string) int {
	var count int

	db := config.InitDB()
	defer db.Close()

	query := `select count(id) from logistics where name = $1`
	db.QueryRow(query, name).Scan(&count)

	return count
}

// get price
func GetMaxPrice(name string) int64 {
	var max_price int64
	db := config.InitDB()
	defer db.Close()

	query := `select max_price from logistics where name = $1`
	db.QueryRow(query, name).Scan(&max_price)

	return max_price
}
