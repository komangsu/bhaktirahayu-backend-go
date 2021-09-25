package models

import (
	"log"
	"restapi/config"
	"restapi/schemas"
)

// Get Max Id
func GetMaxLogisticService() int64 {
	db := config.InitDB()
	defer db.Close()

	query := `select max(id) from logistic_services`

	rows, _ := db.Query(query)
	var max_id int64
	for rows.Next() {
		rows.Scan(&max_id)
	}

	return max_id
}

// Check duplicate name
func CheckLogServiceDuplicateName(name string) int {
	var count int
	db := config.InitDB()
	defer db.Close()

	query := `select count(id) from logistic_services where name = $1`
	db.QueryRow(query, name).Scan(&count)

	return count
}

// Check logistic service by id
func CheckLogServiceById(id int) int {
	var count int

	db := config.InitDB()
	defer db.Close()

	query := `select count(id) from logistic_services where id = $1 `
	db.QueryRow(query, id).Scan(&count)

	return count
}

// Create logistic service
func CreateLogisticService(payload schemas.LogisticServices) int64 {
	db := config.InitDB()
	defer db.Close()

	query := `insert into logistic_services(id,name,margin) values($1,$2,$3) returning id`

	stmt, err := db.Prepare(query)
	if err != nil {
		panic(err)
	}

	log_id := GetMaxLogisticService() + 1

	var last_id int64
	queryErr := stmt.QueryRow(log_id, payload.Name, payload.Margin).Scan(&last_id)
	if queryErr != nil {
		log.Fatal(queryErr)
	}

	return last_id
}

// Update
func UpdateLogisticService(id int, name string, margin int64) {
	db := config.InitDB()
	defer db.Close()

	query := `update logistic_services set name = $1,margin = $2 where id = $3`

	stmt, _ := db.Prepare(query)
	stmt.Exec(name, margin, id)

	return
}

// Delete
func DeleteLogisticService(id int) {
	db := config.InitDB()
	defer db.Close()

	query := `delete from logistic_services where id = $1`

	stmt, _ := db.Prepare(query)
	stmt.Exec(id)

	return
}

// Get All Logistic Services
func GetAllLogisticServices(q string, page schemas.PaginationLogistic) ([]schemas.ListAllLogisticServices, int) {
	logistic_services := []schemas.ListAllLogisticServices{}

	db := config.InitDB()
	defer db.Close()

	query := `select * from logistic_services where name ilike '%' || $1 || '%' order by id desc limit $2 offset $3`

	query_total := `select count(id) from logistic_services`

	var total int
	db.QueryRow(query_total).Scan(&total)

	offset := (page.Page - 1) * page.PerPage
	rows, _ := db.Query(query, q, page.PerPage, offset)
	for rows.Next() {
		var log_services schemas.ListAllLogisticServices

		rows.Scan(&log_services.Id, &log_services.Name)
		logistic_services = append(logistic_services, log_services)
	}

	return logistic_services, total
}

// Get Multiple Logistic Services
func GetMultipleLogisticServices() []schemas.ListAllLogisticServices {
	logistic_services := []schemas.ListAllLogisticServices{}

	db := config.InitDB()
	defer db.Close()

	query := `select * from logistic_services`
	rows, _ := db.Query(query)
	for rows.Next() {
		var log_services schemas.ListAllLogisticServices

		rows.Scan(&log_services.Id, &log_services.Name, &log_services.Margin)
		logistic_services = append(logistic_services, log_services)
	}
	return logistic_services
}
