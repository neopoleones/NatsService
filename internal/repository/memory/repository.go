package memory

import (
	"L0_task/internal/models"
	"L0_task/internal/repository"
	"context"
	"github.com/hashicorp/go-memdb"
	"sync"
)

const (
	tableOrder = "order"
)

var (
	schema *memdb.DBSchema
	once   sync.Once
)

type OrdersRepository struct {
	db *memdb.MemDB
}

func (o OrdersRepository) GetOrders(_ context.Context) ([]*models.Order, error) {
	orders := make([]*models.Order, 0)

	tx := o.db.Txn(false)
	defer tx.Abort()

	cur, err := tx.Get(tableOrder, "id")
	if err != nil {
		return nil, err
	}

	for obj := cur.Next(); obj != nil; obj = cur.Next() {
		if order, ok := obj.(*models.Order); ok {
			orders = append(orders, order)
		}
	}

	return orders, nil
}

func (o OrdersRepository) GetOrderByUID(_ context.Context, uid string) (*models.Order, error) {
	tx := o.db.Txn(false)
	defer tx.Abort()

	obj, err := tx.First(tableOrder, "id", uid)
	if err != nil {
		return nil, err
	}

	if order, ok := obj.(*models.Order); ok {
		return order, nil
	} else {
		return nil, repository.ErrInvalidOrder
	}
}

func (o OrdersRepository) AddOrder(_ context.Context, order *models.Order) error {
	tx := o.db.Txn(true)

	if err := tx.Insert(tableOrder, order); err != nil {
		tx.Abort()
		return err
	}

	tx.Commit()
	return nil
}

func (o OrdersRepository) Close() error {
	// Dummy func (as we don't need to close memdb)
	return nil
}

func NewOrdersRepository() (*OrdersRepository, error) {
	// Sets up the schema for database only once
	once.Do(func() {
		schema = &memdb.DBSchema{
			Tables: map[string]*memdb.TableSchema{
				"order": {
					Name: tableOrder,
					Indexes: map[string]*memdb.IndexSchema{
						"id": {
							Name:   "id",
							Unique: true,
							Indexer: &memdb.StringFieldIndex{
								Field: "OrderUID",
							},
						},
					},
				},
			},
		}
	})

	db, err := memdb.NewMemDB(schema)
	if err != nil {
		return nil, err
	}
	return &OrdersRepository{
		db: db,
	}, nil
}
