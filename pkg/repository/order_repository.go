package repository

import (
	"github.com/jmoiron/sqlx"
	wb_l0 "wb-l0"
)

type OrderPostgres struct {
	db *sqlx.DB
}

func NewOrderPostgres(db *sqlx.DB) *OrderPostgres {
	return &OrderPostgres{db: db}
}

func (o *OrderPostgres) GetOrder(id string) (wb_l0.Order, error) {
	return wb_l0.Order{}, nil
}
