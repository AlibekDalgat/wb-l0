package service

import (
	wb_l0 "wb-l0"
	"wb-l0/pkg/repository"
)

type Order interface {
	GetOrder(id string) (wb_l0.Order, error)
}

type Service struct {
	Order
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Order: NewOrderService(repo),
	}
}
