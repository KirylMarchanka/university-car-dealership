package employee

import (
	"car_dealership/internal/hash"
	"errors"
	"log"
	"time"
)

type Employee struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
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

	date := time.Now()
	id := store(name, email, hashedPass, date)
	if id == 0 {
		return nil
	}

	return &Employee{
		Id:        id,
		Name:      name,
		Email:     email,
		Password:  hashedPass,
		CreatedAt: date,
		UpdatedAt: date,
	}
}

func Find(email string) *Employee {
	emlp := findByEmail(email)
	if emlp == nil {
		return nil
	}

	return emlp
}
