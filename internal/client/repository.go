package client

import (
	"database/sql"
	"errors"
	"log"
)

func findByPhone(phone string) (*Client, error) {
	db := connect()
	defer db.Close()

	var c Client

	err := db.QueryRow("SELECT id, name, phone FROM clients WHERE phone = ?", phone).Scan(
		&c.Id, &c.Name, &c.Phone,
	)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Print(err.Error())
		}

		return nil, err
	}

	return &c, nil
}

func getClients() ([]WithOrderCount, error) {
	db := connect()
	defer db.Close()

	// Query to get clients with orders count
	query := `
		SELECT c.id, c.name, c.phone, COUNT(o.id) AS order_count
		FROM clients c
		LEFT JOIN orders o ON c.id = o.client_id
		GROUP BY c.id
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}
	defer rows.Close()

	var clients []WithOrderCount
	for rows.Next() {
		var client WithOrderCount
		err := rows.Scan(&client.Id, &client.Name, &client.Phone, &client.OrderCount)
		if err != nil {
			log.Print(err.Error())
			return nil, err
		}
		clients = append(clients, client)
	}

	return clients, nil
}

func getClient(id int64) ([]WithOrder, error) {
	db := connect()
	defer db.Close()

	query := `SELECT
    c.id AS client_id,
    c.name AS client_name,
    c.phone AS client_phone,
    o.order_date,
    c2.name AS car_name,
    s.sale_date,
    s.price
FROM clients AS c
         INNER JOIN orders AS o ON c.id = o.client_id
         INNER JOIN cars AS c2 ON o.car_id = c2.id
         INNER JOIN sales AS s ON o.id = s.order_id AND c2.id = s.car_id
WHERE c.id = ?;`

	rows, err := db.Query(query, id)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	defer rows.Close()

	var clientsData []WithOrder

	// Iterate through the query results and populate the struct
	for rows.Next() {
		var clientData WithOrder
		err := rows.Scan(
			&clientData.Id,
			&clientData.Name,
			&clientData.Phone,
			&clientData.OrderDate,
			&clientData.CarName,
			&clientData.SaleDate,
			&clientData.Price,
		)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		clientsData = append(clientsData, clientData)
	}

	return clientsData, nil
}

func destroy(id int64) error {
	db := connect()
	defer db.Close()

	_, err := db.Exec("DELETE FROM clients WHERE id = ?", id)
	if err != nil {
		log.Print(err.Error())
		return err
	}

	return nil
}

func storeClient(name, phone, password string) int64 {
	db := connect()
	defer db.Close()

	ins, err := db.Prepare("INSERT INTO clients (name, phone, password) VALUES (?, ?, ?)")
	if err != nil {
		log.Printf("Unable to insert Client, err: [%s]\n", err.Error())
		return 0
	}

	res, err := ins.Exec(name, phone, password)
	if err != nil {
		log.Printf("Unable to insert Client, err: [%s]\n", err.Error())
		return 0
	}

	id, _ := res.LastInsertId()

	return id
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
