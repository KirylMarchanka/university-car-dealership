package main

import (
	"car_dealership/internal/employee"
	"fmt"
)

func main() {
	err := employee.Validate("John Doe", "jd@gmail.com", "test321")
	if err != nil {
		fmt.Printf("Failed to validate Employee, err: %s", err.Error())

		return
	}

	empl := employee.New("John Doe", "jd@gmail.com", "test321")

	fmt.Println(empl)
}
