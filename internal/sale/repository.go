package sale

import (
	"database/sql"
	"errors"
	"log"
)

func Exists(employeeId, orderId int64) bool {
	db := connect()
	defer db.Close()

	var id int64
	err := db.QueryRow("SELECT (1) FROM sales WHERE employee_id = ? and order_id = ?", employeeId, orderId).Scan(&id)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Printf("Unable to get Sale, err: [%s]\n", err.Error())
		}

		return false
	}

	return true
}

func connect() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "car-dealership"
	dbPass := "password"
	dbName := "car-dealership"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}

	return db
}
