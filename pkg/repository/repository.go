package repository

import (
	"github.com/jmoiron/sqlx"
	wb_l0 "wb-l0"
)

type Order interface {
	AddOrder(orderUID string, data []byte) error
	PullAllOrders() (map[string]wb_l0.Order, error)
}

type Repository struct {
	Order
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Order: NewOrderPostgres(db),
	}
}
