package employee

import (
	"car_dealership/internal/hash"
	"errors"
	"log"
)

type Employee struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type EmployeeSales struct {
	Name       string  `json:"name"`
	Email      string  `json:"email"`
	Sales      []sale  `json:"sales"`
	MinPrice   float64 `json:"min_price"`
	MaxPrice   float64 `json:"max_price"`
	TotalPrice float64 `json:"total_price"`
}

type sale struct {
	CarName  string  `json:"car_name"`
	Price    float64 `json:"price"`
	SaleDate string  `json:"sale_date"`
}

func Validate(name, email, password string) error {
	if len(name) <= 2 || len(name) > 255 {
		return errors.New("incorrect employee name length")
	}

	if len(email) <= 2 || len(email) > 254 {
		return errors.New("incorrect employee email length")
	}

	if existsByEmail(email) {
		return errors.New("non unique employee email")
	}

	if len(password) < 8 {
		return errors.New("incorrect employee password length")
	}

	return nil
}

func New(name, email, password string) *Employee {
	hashedPass, err := hash.Hash(password)
	if err != nil {
		log.Printf("Unable to hash Employee password, err: %s", err.Error())
		return nil
	}

	id := store(name, email, hashedPass)
	if id == 0 {
		return nil
	}

	return &Employee{
		Id:       id,
		Name:     name,
		Email:    email,
		Password: hashedPass,
	}
}

func Find(email string) *Employee {
	emlp := findByEmail(email)
	if emlp == nil {
		return nil
	}

	return emlp
}

func FindById(id int64) *Employee {
	emlp := findById(id)
	if emlp == nil {
		return nil
	}

	return emlp
}

func GetEmployers(minDate, maxDate string, minPrice, maxPrice int64) ([]map[string]interface{}, error) {
	return getEmployers(minDate, maxDate, minPrice, maxPrice)
}

func GetEmployee(id int64) *EmployeeSales {
	employeeSales, err := getWithSales(id)
	if err != nil {
		return nil
	}

	return employeeSales
}

func Delete(id int64) bool {
	db := connect()
	defer db.Close()

	// Delete the car from the database.
	_, err := db.Exec("DELETE FROM employees WHERE id = ?", id)
	if err != nil {
		log.Print(err.Error())
		return false
	}

	return true
}
