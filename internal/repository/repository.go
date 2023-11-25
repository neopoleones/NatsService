package repository

import (
	"L0_task/internal/models"
	"context"
	"errors"
	"io"
)

var (
	ErrInvalidOrder  = errors.New("invalid order")
	ErrAlreadyExists = errors.New("already exists")
)

type OrdersRepository interface {
	io.Closer

	GetOrders(context.Context) ([]*models.Order, error)
	GetOrderByUID(context.Context, string) (*models.Order, error)
	AddOrder(context.Context, *models.Order) error
}
