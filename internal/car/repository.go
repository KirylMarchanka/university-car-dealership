package car

import (
	"database/sql"
	"errors"
	"log"
	"strings"
)

func store(
	manufacturerId int64,
	name string,
	fuel string,
	fuelCapacity float32,
	engine string,
	enginePower float32,
	engineCapacity float32,
	maxSpeed int32,
	acceleration float32,
) int64 {
	db := connect()
	defer db.Close()

	ins, err := db.Prepare("INSERT INTO cars (manufacturer_id, name, fuel, fuel_capacity, engine, engine_power, engine_capacity, max_speed, acceleration) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Printf("Unable to insert Car, err: [%s]\n", err.Error())
		return 0
	}

	res, err := ins.Exec(manufacturerId, name, fuel, fuelCapacity, engine, enginePower, engineCapacity, maxSpeed, acceleration)
	if err != nil {
		log.Printf("Unable to insert Car, err: [%s]\n", err.Error())
		return 0
	}

	id, _ := res.LastInsertId()

	return id
}

func (c *Car) update() error {
	if !carExists(c.Id) {
		return errors.New("car ID not found")
	}

	db := connect()
	defer db.Close()

	updateSQL := `
		UPDATE cars
		SET manufacturer_id = ?, name = ?, fuel = ?, fuel_capacity = ?, engine = ?, engine_power = ?, engine_capacity = ?, max_speed = ?, acceleration = ?
		WHERE id = ?
	`
	_, err := db.Exec(updateSQL, c.ManufacturerId, c.Name, c.Fuel, c.FuelCapacity, c.Engine, c.EnginePower, c.EngineCapacity, c.MaxSpeed, c.Acceleration, c.Id)
	if err != nil {
		return err
	}

	return nil
}

func (c *Car) delete() error {
	if !carExists(c.Id) {
		return errors.New("car ID not found")
	}

	db := connect()
	defer db.Close()

	// Delete the car from the database.
	_, err := db.Exec("DELETE FROM cars WHERE id = ?", c.Id)
	if err != nil {
		return err
	}

	return nil
}

func carExists(id int64) bool {
	db := connect()
	defer db.Close()

	var existingCarID int64

	err := db.QueryRow("SELECT id FROM cars WHERE id = ?", id).Scan(&existingCarID)
	if errors.Is(err, sql.ErrNoRows) {
		return false
	} else if err != nil {
		log.Fatal(err)
	}

	return true
}

func selectCars(manufacturerID int64, name, fuel, orderBy, orderDirection string) ([]Car, error) {
	// Construct the base SELECT query.
	selectSQL := `
		SELECT c.id, c.manufacturer_id, c.name, c.fuel, c.fuel_capacity, c.engine, c.engine_power, c.engine_capacity, c.max_speed, c.acceleration, m.id, m.name AS manufacturer_name
		FROM cars c
		LEFT JOIN manufacturers m ON c.manufacturer_id = m.id
		WHERE 1=1
	`

	// Add filters based on provided criteria.
	var args []interface{}

	if manufacturerID != 0 {
		selectSQL += " AND manufacturer_id = ?"
		args = append(args, manufacturerID)
	}
	if name != "" {
		selectSQL += " AND name = ?"
		args = append(args, name)
	}
	if fuel != "" {
		selectSQL += " AND fuel = ?"
		args = append(args, fuel)
	}

	db := connect()
	defer db.Close()

	// Add ordering criteria.
	orderBy = strings.ToLower(orderBy)
	if orderBy == "max_speed" || orderBy == "acceleration" {
		selectSQL += " ORDER BY " + orderBy
		if orderDirection == "desc" {
			selectSQL += " DESC"
		} else {
			selectSQL += " ASC"
		}
	}

	rows, err := db.Query(selectSQL, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cars []Car
	for rows.Next() {
		var car Car
		err := rows.Scan(
			&car.Id, &car.ManufacturerId, &car.Name, &car.Fuel, &car.FuelCapacity, &car.Engine, &car.EnginePower,
			&car.EngineCapacity, &car.MaxSpeed, &car.Acceleration, &car.Manufacturer.Id, &car.Manufacturer.Name,
		)
		if err != nil {
			return nil, err
		}

		cars = append(cars, car)
	}

	return cars, nil
}

func findById(id int64) (Car, error) {
	selectSQL := `
		SELECT c.id, c.manufacturer_id, c.name, c.fuel, c.fuel_capacity, c.engine, c.engine_power, c.engine_capacity, c.max_speed, c.acceleration, m.id, m.name AS manufacturer_name
		FROM cars c
		LEFT JOIN manufacturers m ON c.manufacturer_id = m.id
		WHERE c.id = ?
	`

	db := connect()
	defer db.Close()

	var car Car

	err := db.QueryRow(selectSQL, id).Scan(
		&car.Id, &car.ManufacturerId, &car.Name, &car.Fuel, &car.FuelCapacity, &car.Engine, &car.EnginePower,
		&car.EngineCapacity, &car.MaxSpeed, &car.Acceleration, &car.Manufacturer.Id, &car.Manufacturer.Name,
	)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Printf("Unable to get Car by id, err: [%s]\n", err.Error())
		}

		return Car{}, err
	}

	return car, nil
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
