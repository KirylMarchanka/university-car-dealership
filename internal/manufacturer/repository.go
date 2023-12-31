package manufacturer

import (
	"database/sql"
	"errors"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func store(name string) int64 {
	db := connect()
	defer db.Close()

	ins, err := db.Prepare("INSERT INTO manufacturers (name) VALUES (?)")
	if err != nil {
		log.Printf("Unable to insert Manufacturer, err: [%s]\n", err.Error())
		return 0
	}

	res, err := ins.Exec(name)
	if err != nil {
		log.Printf("Unable to insert Manufacturer, err: [%s]\n", err.Error())
		return 0
	}

	id, _ := res.LastInsertId()

	return id
}

func getById(id int64) (int64, string, error) {
	db := connect()
	defer db.Close()

	var name string

	err := db.QueryRow("SELECT * FROM manufacturers WHERE id = ?", id).Scan(&id, &name)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Printf("Unable to get Manufacturer by ID, err: [%s]\n", err.Error())
		}

		return 0, "", err
	}

	return id, name, nil
}

func existsByName(name string) bool {
	db := connect()
	defer db.Close()

	err := db.QueryRow("SELECT (1) FROM manufacturers WHERE name = ?", name).Scan(&name)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Printf("Unable to get Manufacturer by name, err: [%s]\n", err.Error())
		}

		return false
	}

	return true
}

func ExistsById(id int64) bool {
	db := connect()
	defer db.Close()

	err := db.QueryRow("SELECT id FROM manufacturers WHERE id = ?", id).Scan(&id)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Printf("Unable to get Manufacturer by ID, err: [%s]\n", err.Error())
		}

		return false
	}

	return true
}

func get() ([]Manufacturer, error) {
	db := connect()
	defer db.Close()

	rows, err := db.Query("SELECT id, name FROM manufacturers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var manufacturers []Manufacturer

	for rows.Next() {
		var manufacturer Manufacturer
		if err := rows.Scan(&manufacturer.Id, &manufacturer.Name); err != nil {
			return nil, err
		}
		manufacturers = append(manufacturers, manufacturer)
	}

	return manufacturers, nil
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
