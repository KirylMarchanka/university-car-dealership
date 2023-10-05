package employee

import (
	"database/sql"
	"errors"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func store(name, email, password string, date time.Time) int64 {
	db := connect()
	defer db.Close()

	ins, err := db.Prepare("INSERT INTO employees (name, email, password, created_at, updated_at) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Printf("Unable to insert Employee, err: [%s]\n", err.Error())
		return 0
	}

	formattedDate := date.Format("2006-01-02 15:04:05")
	res, err := ins.Exec(name, email, password, formattedDate, formattedDate)
	if err != nil {
		log.Printf("Unable to insert Employee, err: [%s]\n", err.Error())
		return 0
	}

	id, _ := res.LastInsertId()

	return id
}

func existsByEmail(email string) bool {
	db := connect()
	defer db.Close()

	err := db.QueryRow("SELECT (1) FROM employees WHERE email = ?", email).Scan(&email)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Printf("Unable to get Employee by email, err: [%s]\n", err.Error())
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
