package postgres

const (
	queryGetItemsByOrderID = `
		SELECT
	    	chrt_id, track_number,
	    	price, rid, name, sale,
	    	size, total_price, nm_id,
	    	brand, status
		FROM wb_item
		WHERE order_id=$1
	`
	queryGetDeliveryByID = `
		SELECT
			name, phone, zip,
			city, address, region, email
		FROM wb_delivery
		WHERE delivery_id=$1
		LIMIT 1
	`
	queryGetPaymentByID = `
		SELECT
			transaction, request_id, currency,
			provider, amount, payment_dt,
			bank, delivery_cost, goods_total,
			custom_fee
		FROM wb_payment
		WHERE payment_id=$1
		LIMIT 1
	`
	queryGetOrders = `
		SELECT
			order_id, delivery_id, payment_id,
			shard_key, sm_id, oof_shard, date_created,
			internal_signature, raw_order_id, track_number,
			customer_id, delivery_service, entry, locale
		FROM wb_order
	`
	queryGetOrderByRawID = `
		SELECT
			order_id, delivery_id, payment_id,
			shard_key, sm_id, oof_shard, date_created,
			internal_signature, raw_order_id, track_number,
			customer_id, delivery_service, entry, locale
		FROM wb_order
		WHERE raw_order_id=$1
		LIMIT 1
	`

	queryInsertOrder = `
		INSERT INTO wb_order(delivery_id, payment_id, shard_key, sm_id, oof_shard, date_created, internal_signature, raw_order_id, track_number, customer_id, delivery_service, entry, locale)
		VALUES (@delivery_id, @payment_id, @shard_key, @sm_id, @oof_shard, @date_created, @internal_signature, @raw_order_id, @track_number, @customer_id, @delivery_service, @entry, @locale)
		RETURNING order_id
		
	`
	queryInsertItem = `
		INSERT INTO wb_item(order_id, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
		VALUES (@order_id, @chrt_id, @track_number, @price, @rid, @name, @sale, @size, @total_price, @nm_id, @brand, @status)
	`
	queryInsertPayment = `
		INSERT INTO wb_payment(transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
		VALUES (@transaction, @request_id, @currency, @provider, @amount, @payment_dt, @bank, @delivery_cost, @goods_total, @custom_fee)
		RETURNING payment_id
	`
	queryInsertDelivery = `
		INSERT INTO wb_delivery(name, phone, zip, city, address, region, email)
		VALUES (@name, @phone, @zip, @city, @address, @region, @email)
		RETURNING delivery_id
	`
)
