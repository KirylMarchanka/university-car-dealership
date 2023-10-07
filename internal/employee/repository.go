package employee

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strconv"
)

func store(name, email, password string) int64 {
	db := connect()
	defer db.Close()

	ins, err := db.Prepare("INSERT INTO employees (name, email, password) VALUES (?, ?, ?)")
	if err != nil {
		log.Printf("Unable to insert Employee, err: [%s]\n", err.Error())
		return 0
	}

	res, err := ins.Exec(name, email, password)
	if err != nil {
		log.Printf("Unable to insert Employee, err: [%s]\n", err.Error())
		return 0
	}

	id, _ := res.LastInsertId()

	return id
}

func findByEmail(email string) *Employee {
	db := connect()
	defer db.Close()

	var empl Employee

	err := db.QueryRow("SELECT * FROM employees WHERE email = ?", email).Scan(
		&empl.Id, &empl.Name, &empl.Email, &empl.Password,
	)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Printf("Unable to get Employee by email, err: [%s]\n", err.Error())
		}

		return nil
	}

	return &empl
}

func findById(id int64) *Employee {
	db := connect()
	defer db.Close()

	var empl Employee

	err := db.QueryRow("SELECT * FROM employees WHERE id = ?", id).Scan(
		&empl.Id, &empl.Name, &empl.Email, &empl.Password,
	)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Printf("Unable to get Employee by id, err: [%s]\n", err.Error())
		}

		return nil
	}

	return &empl
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

func getEmployers(minDate, maxDate string, minPrice, maxPrice int64) ([]map[string]interface{}, error) {
	db := connect()
	defer db.Close()

	// Execute the SQL query
	query := `
    SELECT e.name, e.email, COUNT(s.id) as sales_count, COALESCE(MAX(s.price), 0) as sales_max_price, COALESCE(AVG(s.price), 0.00) as sales_avg_price
    FROM employees e
    LEFT JOIN sales s ON e.id = s.employee_id
    `

	if minDate != "" || maxDate != "" || minPrice != 0 || maxPrice != 0 {
		query += "WHERE "
	}

	var hasWhere bool

	if minPrice != 0 {
		query += generateFilterClause("s.price", ">= ", strconv.FormatInt(minPrice, 10), hasWhere)
		hasWhere = true
	}

	if maxPrice != 0 {
		query += generateFilterClause("s.price", "<= ", strconv.FormatInt(maxPrice, 10), hasWhere)
		hasWhere = true
	}

	if minDate != "" {
		query += generateFilterClause("s.sale_date", ">= ", minDate, hasWhere)
		hasWhere = true
	}

	if maxDate != "" {
		query += generateFilterClause("s.sale_date", "<= ", maxDate, hasWhere)
		hasWhere = true
	}

	query += ` GROUP BY e.name, e.email`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Create a slice to hold the result data
	var employees []map[string]interface{}

	// Iterate through the rows and populate the result slice
	for rows.Next() {
		var name, email string
		var salesCount, salesMaxPrice int
		var salesAvgPrice float32
		if err := rows.Scan(&name, &email, &salesCount, &salesMaxPrice, &salesAvgPrice); err != nil {
			return nil, err
		}

		employeeData := map[string]interface{}{
			"name":            name,
			"email":           email,
			"sales_count":     salesCount,
			"sales_max_price": salesMaxPrice,
			"sales_avg_price": salesAvgPrice,
		}
		employees = append(employees, employeeData)
	}

	return employees, nil
}

func getWithSales(id int64) (*EmployeeSales, error) {
	db := connect()
	defer db.Close()

	query := `
            SELECT e.name, e.email, s.price, s.sale_date, c.name AS car_name
            FROM employees e
            LEFT JOIN sales s ON e.id = s.employee_id
            LEFT JOIN cars c ON s.car_id = c.id
            WHERE e.id = ?`

	rows, err := db.Query(query, id)
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}
	defer rows.Close()

	employeeSale := EmployeeSales{}
	minPrice := float64(0)
	maxPrice := float64(0)
	totalPrice := float64(0)

	for rows.Next() {
		var name, email string
		var nullableCarName, nullableSaleDate sql.NullString
		var nullablePrice sql.NullFloat64
		err := rows.Scan(&name, &email, &nullablePrice, &nullableSaleDate, &nullableCarName)
		if err != nil {
			log.Print(err.Error())
			return nil, err
		}

		employeeSale.Name = name
		employeeSale.Email = email

		var price float64
		if !nullablePrice.Valid {
			continue
		}

		price = nullablePrice.Float64

		if minPrice == 0 || price < minPrice {
			minPrice = price
		}
		if price > maxPrice {
			maxPrice = price
		}
		totalPrice += price

		employeeSale.Sales = append(employeeSale.Sales, sale{CarName: nullableCarName.String, Price: price, SaleDate: nullableSaleDate.String})
	}

	employeeSale.MinPrice = minPrice
	employeeSale.MaxPrice = maxPrice
	employeeSale.TotalPrice = totalPrice

	return &employeeSale, nil
}

// Helper function to generate filter conditions based on provided values
func generateFilterClause(column, condition, value string, hasWhere bool) string {
	var query string
	if hasWhere {
		query += " AND "
	}

	if value != "" {
		return query + fmt.Sprintf("%s %s %s", column, condition, value)
	}

	return query + "1=1" // True condition when value is empty (filter not provided)
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
