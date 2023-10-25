package service

import (
	wb_l0 "wb-l0"
	"wb-l0/pkg/repository"
)

type OrderService struct {
	repo *repository.Repository
}

func NewOrderService(repo *repository.Repository) *OrderService {
	return &OrderService{repo: repo}
}

func (o *OrderService) GetOrder(id string) (wb_l0.Order, error) {
	return o.repo.GetOrder(id)
}

func (o *OrderService) AddOrder(orderUID string, data []byte) error {
	return o.repo.AddOrder(orderUID, data)
}
