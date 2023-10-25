package repository

import (
	"errors"
	"fmt"
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

func (o *OrderPostgres) AddOrder(orderUID string, data []byte) error {
	var count int
	query := fmt.Sprintf("SELECT COUNT(id) FROM %s WHERE id = $1", itemsTable)
	err := o.db.QueryRow(query, orderUID).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("Заказ с таким именем уже существует.")
	}
	query = fmt.Sprintf("INSERT INTO %s (id, description) values ($1, $2)", itemsTable)
	_, err = o.db.Exec(query, orderUID, data)
	return err
}
