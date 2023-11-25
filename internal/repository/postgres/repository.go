package postgres

import (
	"L0_task/internal/config"
	"L0_task/internal/models"
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type scanner interface {
	Scan(dest ...any) error
}

type OrdersRepository struct {
	c *pgxpool.Pool
}

type orderExtension struct {
	models.Order

	orderID    int
	deliveryID int
	paymentID  int
}

func (o OrdersRepository) enrichOrder(ctx context.Context, ext *orderExtension) error {
	if err := o.setItemsByOrderID(ctx, ext); err != nil {
		return err
	}

	if err := o.setPaymentByID(ctx, ext); err != nil {
		return err
	}

	return o.setDeliveryByID(ctx, ext)
}

func (o OrdersRepository) setItemsByOrderID(ctx context.Context, ext *orderExtension) error {
	allItems := make([]*models.Item, 0)

	rows, err := o.c.Query(ctx, queryGetItemsByOrderID, ext.orderID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.Item

		err := rows.Scan(
			&item.ChrtId, &item.TrackNumber, &item.Price,
			&item.Rid, &item.Name, &item.Sale, &item.Size,
			&item.TotalPrice, &item.NmId, &item.Brand, &item.Status,
		)

		if err != nil {
			return err
		}
		allItems = append(allItems, &item)
	}

	ext.Items = allItems
	return nil
}

func (o OrdersRepository) setDeliveryByID(ctx context.Context, ext *orderExtension) error {
	var delivery models.Delivery

	err := o.c.QueryRow(ctx, queryGetDeliveryByID, ext.deliveryID).Scan(
		&delivery.Name, &delivery.Phone, &delivery.Zip,
		&delivery.City, &delivery.Address, &delivery.Region,
		&delivery.Email,
	)
	if err != nil {
		return err
	}

	ext.Delivery = &delivery
	return nil
}

func (o OrdersRepository) setPaymentByID(ctx context.Context, ext *orderExtension) error {
	var payment models.Payment

	err := o.c.QueryRow(ctx, queryGetPaymentByID, ext.paymentID).Scan(
		&payment.Transaction, &payment.RequestId, &payment.Currency,
		&payment.Provider, &payment.Amount, &payment.PaymentDt,
		&payment.Bank, &payment.DeliveryCost, &payment.GoodsTotal,
		&payment.CustomFee,
	)
	if err != nil {
		return err
	}

	ext.Payment = &payment
	return nil
}

func (o OrdersRepository) rowToOrder(ctx context.Context, row scanner) (*models.Order, error) {
	var ext orderExtension

	err := row.Scan(
		&ext.orderID, &ext.deliveryID, &ext.paymentID,
		&ext.ShardKey, &ext.SmId, &ext.OofShard, &ext.DateCreated,
		&ext.InternalSignature, &ext.OrderUID, &ext.TrackNumber,
		&ext.CustomerId, &ext.DeliveryService, &ext.Entry, &ext.Locale,
	)

	if err != nil {
		return nil, err
	}
	if err = o.enrichOrder(ctx, &ext); err != nil {
		return nil, err
	}

	return &ext.Order, nil
}

func (o OrdersRepository) GetOrders(ctx context.Context) ([]*models.Order, error) {
	orders := make([]*models.Order, 0)

	rows, err := o.c.Query(ctx, queryGetOrders)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		order, err := o.rowToOrder(ctx, rows)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (o OrdersRepository) GetOrderByUID(ctx context.Context, rawID string) (*models.Order, error) {
	return o.rowToOrder(ctx, o.c.QueryRow(ctx, queryGetOrderByRawID, rawID))
}

func (o OrdersRepository) AddOrder(ctx context.Context, order *models.Order) error {
	ext := orderExtension{Order: *order}

	// Insert delivery
	err := o.c.QueryRow(ctx, queryInsertDelivery, pgx.NamedArgs{
		"name":    ext.Delivery.Name,
		"phone":   ext.Delivery.Phone,
		"zip":     ext.Delivery.Zip,
		"city":    ext.Delivery.City,
		"address": ext.Delivery.Address,
		"region":  ext.Delivery.Region,
		"email":   ext.Delivery.Email,
	}).Scan(&ext.deliveryID)
	if err != nil {
		return err
	}

	// Insert payment
	err = o.c.QueryRow(ctx, queryInsertPayment, pgx.NamedArgs{
		"transaction":   ext.Payment.Transaction,
		"request_id":    ext.Payment.RequestId,
		"currency":      ext.Payment.Currency,
		"provider":      ext.Payment.Provider,
		"amount":        ext.Payment.Amount,
		"payment_dt":    ext.Payment.PaymentDt,
		"bank":          ext.Payment.Bank,
		"delivery_cost": ext.Payment.DeliveryCost,
		"goods_total":   ext.Payment.GoodsTotal,
		"custom_fee":    ext.Payment.CustomFee,
	}).Scan(&ext.paymentID)
	if err != nil {
		return err
	}

	// Insert Order
	err = o.c.QueryRow(ctx, queryInsertOrder, pgx.NamedArgs{
		"delivery_id":        ext.deliveryID,
		"payment_id":         ext.paymentID,
		"shard_key":          ext.Order.ShardKey,
		"sm_id":              ext.Order.SmId,
		"oof_shard":          ext.Order.OofShard,
		"date_created":       ext.Order.DateCreated,
		"internal_signature": ext.Order.InternalSignature,
		"raw_order_id":       ext.Order.OrderUID,
		"track_number":       ext.Order.TrackNumber,
		"customer_id":        ext.Order.CustomerId,
		"delivery_service":   ext.Order.DeliveryService,
		"entry":              ext.Order.Entry,
		"locale":             ext.Order.Locale,
	}).Scan(&ext.orderID)
	if err != nil {
		return err
	}

	for _, item := range ext.Items {
		_, err := o.c.Exec(ctx, queryInsertItem, pgx.NamedArgs{
			"order_id":     ext.orderID,
			"chrt_id":      item.ChrtId,
			"track_number": item.TrackNumber,
			"price":        item.Price,
			"rid":          item.Rid,
			"name":         item.Name,
			"sale":         item.Sale,
			"size":         item.Size,
			"total_price":  item.TotalPrice,
			"nm_id":        item.NmId,
			"brand":        item.Brand,
			"status":       item.Status,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (o OrdersRepository) Close() error {
	o.c.Close()
	return nil
}

func NewOrdersRepository(ctx context.Context, cfg *config.Postgres) (*OrdersRepository, error) {
	conn, err := pgxpool.New(ctx, cfg.AsSchema())
	if err != nil {
		return nil, err
	}

	return &OrdersRepository{
		c: conn,
	}, nil
}
