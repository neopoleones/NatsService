package dual

import (
	"L0_task/internal/config"
	"L0_task/internal/models"
	"L0_task/internal/repository"
	"L0_task/internal/repository/memory"
	"L0_task/internal/repository/postgres"
	"context"
)

type OrdersRepository struct {
	memoryComponent *memory.OrdersRepository
	pgComponent     *postgres.OrdersRepository
}

func (o OrdersRepository) Close() error {
	_ = o.memoryComponent.Close()
	return o.pgComponent.Close()
}

func (o OrdersRepository) GetOrders(ctx context.Context) ([]*models.Order, error) {
	return o.memoryComponent.GetOrders(ctx)
}

func (o OrdersRepository) GetOrderByUID(ctx context.Context, uid string) (*models.Order, error) {
	order, err := o.memoryComponent.GetOrderByUID(ctx, uid)

	// If we found our order in the memdb, then we'll have it in postgres instance. Returning the value right here
	if err == nil {
		return order, nil
	}

	order, _ = o.pgComponent.GetOrderByUID(ctx, uid)
	if order != nil {
		// We have order in pgsql, but not in memdb. Adding it as cache
		_ = o.memoryComponent.AddOrder(ctx, order)
	}
	return order, err
}

func (o OrdersRepository) AddOrder(ctx context.Context, order *models.Order) error {
	if v, _ := o.memoryComponent.GetOrderByUID(ctx, order.OrderUID); v != nil {
		return repository.ErrAlreadyExists
	}

	if err := o.memoryComponent.AddOrder(ctx, order); err != nil {
		return err
	}
	return o.pgComponent.AddOrder(ctx, order)
}

func NewOrdersRepository(ctx context.Context, cfg *config.Postgres) (*OrdersRepository, error) {
	pgc, err := postgres.NewOrdersRepository(ctx, cfg)
	if err != nil {
		return nil, err
	}

	memc, err := memory.NewOrdersRepository()
	if err != nil {
		return nil, err
	}

	// Update state of memdb cache part
	orders, err := pgc.GetOrders(ctx)
	if err != nil {
		return nil, err
	}

	for _, order := range orders {
		if err := memc.AddOrder(ctx, order); err != nil {
			return nil, err
		}
	}

	return &OrdersRepository{
		memc, pgc,
	}, nil
}
