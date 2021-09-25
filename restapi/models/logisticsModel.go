package models

import (
	"log"
	"restapi/config"
	"restapi/schemas"
	"time"
)

func GetMaxIdLogistic() int64 {
	db := config.InitDB()
	defer db.Close()

	query := `select max(id) from logistics`

	rows, _ := db.Query(query)
	var max_id int64
	for rows.Next() {
		rows.Scan(&max_id)
	}
	return max_id
}

// check logistic name duplicate
func CheckLogisticDuplicateName(name string) int {
	var count int

	db := config.InitDB()
	defer db.Close()

	query := `select count(id) from logistics where name = $1`
	db.QueryRow(query, name).Scan(&count)

	return count
}

// check logistic by id
func CheckLogisticById(id int) int {
	var count int

	db := config.InitDB()
	defer db.Close()

	query := `select count(id) from logistics where id = $1`
	db.QueryRow(query, id).Scan(&count)

	return count
}

// check logistic type
func CheckLogisticType(id int64) int {
	var count int

	db := config.InitDB()
	defer db.Close()

	query := `select count(id) from logistic_services where id = $1`
	db.QueryRow(query, id).Scan(&count)

	return count
}

// Create Logistic
func CreateLogistic(logis schemas.LogisticsSchema) {

	db := config.InitDB()
	defer db.Close()

	query := `insert into logistics(id,logistic_type,name,component,max_price,price,expired_date)
values($1,$2,$3,$4,$5,$6,$7) returning id`

	stmt, err := db.Prepare(query)
	if err != nil {
		panic(err)
	}

	log_id := GetMaxIdLogistic() + 1

	// parse time
	layout := "02-01-2006"
	t, err := time.Parse(layout, logis.ExpiredDate)

	var last_id int64
	queryErr := stmt.QueryRow(log_id, logis.LogisticType, logis.Name,
		logis.Component, logis.MaxPrice, logis.Price, t).Scan(&last_id)

	if queryErr != nil {
		log.Fatal(queryErr)
	}

}

// Update Logistic
func UpdateLogistic(id int, logistic schemas.UpdateLogisticsSchema) {
	db := config.InitDB()
	defer db.Close()

	query := `update logistics set logistic_type = $1,name = $2,
	stock = $3,
	component = $4,
	max_price = $5,
	price = $6,
	expired_date = $7,
	updated_at = $8
	where id = $9`
	stmt, err := db.Prepare(query)
	if err != nil {
		panic(err)
	}
	// parse time
	layout := "02-01-2006"
	t, err := time.Parse(layout, logistic.ExpiredDate)

	// add time update
	loc, _ := time.LoadLocation("Asia/Kuala_Lumpur")
	logistic.UpdatedAt = time.Now().In(loc)

	stmt.Exec(logistic.LogisticType, logistic.Name, logistic.Stock, logistic.Component, logistic.MaxPrice,
		logistic.Price, t, logistic.UpdatedAt, id)

	return
}

// Delete Logistic
func DeleteLogistic(id int) {
	db := config.InitDB()
	defer db.Close()

	query := `delete from logistics where id = $1`

	stmt, _ := db.Prepare(query)
	stmt.Exec(id)

	return
}

// Get All Logistic
func GetAllLogistic(q string, page schemas.PaginationLogistic) ([]schemas.ListAllLogistics, int) {
	logistics := []schemas.ListAllLogistics{}

	db := config.InitDB()
	defer db.Close()

	query := `select * from logistics where name ilike '%' || $1 || '%' order by id desc limit $2 offset $3`

	query_total := `select count(*) from logistics`

	var total int
	db.QueryRow(query_total).Scan(&total)

	offset := (page.Page - 1) * page.PerPage
	rows, _ := db.Query(query, q, page.PerPage, offset)
	for rows.Next() {
		var logis schemas.ListAllLogistics

		err_scan := rows.Scan(&logis.Id, &logis.LogisticType, &logis.Name, &logis.Stock, &logis.Component, &logis.MaxPrice, &logis.Price, &logis.ExpiredDate, &logis.CreatedAt, &logis.UpdatedAt)
		if err_scan != nil {
			panic(err_scan)
		}
		logistics = append(logistics, logis)
	}

	return logistics, total
}
