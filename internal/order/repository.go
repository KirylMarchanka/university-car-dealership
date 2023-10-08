package order

import (
	"database/sql"
	"errors"
	"log"
)

func Exists(id int64) bool {
	db := connect()
	defer db.Close()

	err := db.QueryRow("SELECT (1) FROM orders WHERE id = ?", id).Scan(&id)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Printf("Unable to get Order by id, err: [%s]\n", err.Error())
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
