package repository

import (
	"encoding/json"
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

func (o *OrderPostgres) PullAllOrders() (map[string]wb_l0.Order, error) {
	var outputs []wb_l0.Output
	query := fmt.Sprintf("SELECT id, description FROM %s", itemsTable)
	err := o.db.Select(&outputs, query)
	orders := make(map[string]wb_l0.Order)
	for _, output := range outputs {
		var order wb_l0.Order
		err := json.Unmarshal(output.Description, &order)
		if err != nil {
			return orders, errors.New("Получены некорреткные данные из БД:" + err.Error())
		}
		orders[output.OrderUID] = order
	}
	return orders, err
}

func (o *OrderPostgres) AddOrder(orderUID string, data []byte) error {
	query := fmt.Sprintf("INSERT INTO %s (id, description) values ($1, $2)", itemsTable)
	_, err := o.db.Exec(query, orderUID, data)
	return err
}
