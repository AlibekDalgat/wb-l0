package service

import (
	"errors"
	wb_l0 "wb-l0"
	"wb-l0/pkg/repository"
)

type OrderService struct {
	repo  *repository.Repository
	cache map[string]wb_l0.Order
}

func NewOrderService(repo *repository.Repository) *OrderService {
	return &OrderService{
		repo:  repo,
		cache: make(map[string]wb_l0.Order),
	}
}

func (o *OrderService) GetOrder(id string) (wb_l0.Order, error) {
	if _, ok := o.cache[id]; !ok {
		return wb_l0.Order{}, errors.New("Нет заказа с данным id: " + id)
	}
	return o.cache[id], nil
}

func (o *OrderService) AddOrder(order wb_l0.Order, data []byte) error {
	if order.OrderUID == "" {
		return errors.New("Пустой OrderUID")
	}
	if _, ok := o.cache[order.OrderUID]; ok {
		return errors.New("Заказ с таким именем уже существует.")
	}
	o.cache[order.OrderUID] = order
	return o.repo.AddOrder(order.OrderUID, data)
}

func (o *OrderService) PullAllOrders() error {
	var err error
	o.cache, err = o.repo.PullAllOrders()
	return err
}
