package service

import (
	wb_l0 "wb-l0"
	"wb-l0/pkg/repository"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go.
type Order interface {
	GetOrder(id string) (wb_l0.Order, error)
	AddOrder(orderUID wb_l0.Order, data []byte) error
	PullAllOrders() error
}

type Service struct {
	Order
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Order: NewOrderService(repo),
	}
}
