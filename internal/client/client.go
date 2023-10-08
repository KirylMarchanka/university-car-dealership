package client

import (
	"car_dealership/internal/hash"
	"database/sql"
	"errors"
	"unicode/utf8"
)

type Client struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type WithOrderCount struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	Phone      string `json:"phone"`
	OrderCount int64  `json:"order_count"`
}

type WithOrder struct {
	Id        int64   `json:"id"`
	Name      string  `json:"name"`
	Phone     string  `json:"phone"`
	OrderDate string  `json:"order_date"`
	CarName   string  `json:"car_name"`
	SaleDate  string  `json:"sale_date"`
	Price     float64 `json:"price"`
}

func GetByPhone(phone string) *Client {
	c, err := findByPhone(phone)
	if err != nil {
		return nil
	}

	return c
}

func GetClients() *[]WithOrderCount {
	clients, err := getClients()
	if err != nil {
		return nil
	}

	return &clients
}

func GetClient(id int64) []WithOrder {
	client, err := getClient(id)
	if err != nil {
		return nil
	}

	return client
}

func (c *Client) Delete() bool {
	err := destroy(c.Id)
	if err != nil {
		return false
	}

	return true
}

func New(name, phone, password string) *Client {
	password, err := hash.Hash(password)
	if err != nil {
		return nil
	}

	cId := storeClient(name, phone, password)
	if cId == 0 {
		return nil
	}

	return &Client{
		Id:       cId,
		Name:     name,
		Phone:    phone,
		Password: "",
	}
}

func (c *Client) Validate() error {
	if utf8.RuneCountInString(c.Name) < 2 || utf8.RuneCountInString(c.Name) > 255 {
		return errors.New("incorrect name length")
	}

	if _, err := findByPhone(c.Phone); !errors.Is(err, sql.ErrNoRows) {
		return errors.New("phone should be unique")
	}

	if utf8.RuneCountInString(c.Password) < 8 {
		return errors.New("password too short")
	}

	return nil
}
