package repository

import (
	"github.com/jmoiron/sqlx"
	wb_l0 "wb-l0"
)

type Order interface {
	GetOrder(id string) (wb_l0.Order, error)
}

type Repository struct {
	Order
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Order: NewOrderPostgres(db),
	}
}
