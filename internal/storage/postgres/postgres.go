//Пакет, реализующий хранение заказов в БД и их подгрузку в кэщ.

package postgres

import (
	"context"
	"log"
	"wb-lvl0/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Структура БД.
type DB struct {
	db *pgxpool.Pool
}

// Функция конструктор для создания нового экземпляра постгрес и миграции таблиц в БД.
func New(dsn string) (*DB, error) {
	db, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	s := DB{
		db: db,
	}
	log.Println("выполнено подключение к БД")
	return &s, nil
}

// Функция для подгрузки всех заказов.
func (d *DB) GetAll() ([]models.Order, error) {
	var orders []models.Order
	rows, err := d.db.Query(context.Background(), `
	SELECT
	orders.order_uid,
	orders.track_number,
	orders."entry",
	orders.locale,
	orders.internal_signature,
	orders.customer_id,
	orders.delivery_service,
	orders.shardkey,
	orders.sm_id,
	orders.date_created,
	orders.oof_shard,
	delivery."name",
	delivery.phone,
	delivery.zip,
	delivery.city,
	delivery."address",
	delivery.region,
	delivery.email,
	payment."transaction",
	payment.request_id,
	payment.currency,
	payment."provider",
	payment.amount,
	payment.payment_dt,
	payment.bank,
	payment.delivery_cost,
	payment.goods_total,
	payment.custom_fee
	FROM orders, delivery, payment
	WHERE orders.order_uid=delivery.order_ref AND orders.order_uid=payment.order_ref
	ORDER BY orders.order_uid;
	`)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var o models.Order
		err = rows.Scan(
			&o.OrderUid,
			&o.TrackNumber,
			&o.Entry,
			&o.Locale,
			&o.InternalSignature,
			&o.CustomerId,
			&o.DeliveryService,
			&o.ShardKey,
			&o.SmID,
			&o.DateCreated,
			&o.OffShard,
			&o.Delivery.Name,
			&o.Delivery.Phone,
			&o.Delivery.Zip,
			&o.Delivery.City,
			&o.Delivery.Address,
			&o.Delivery.Region,
			&o.Delivery.Email,
			&o.Payment.Transaction,
			&o.Payment.RequestId,
			&o.Payment.Currency,
			&o.Payment.Provider,
			&o.Payment.Amount,
			&o.Payment.PaymentDt,
			&o.Payment.Bank,
			&o.Payment.DeliveryCost,
			&o.Payment.GoodsTotal,
			&o.Payment.CustomFee,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, o)

		if rows.Err() != nil {
			return nil, rows.Err()
		}
	}
	for v := range orders {
		uid := orders[v].OrderUid
		rows, err := d.db.Query(context.Background(), `
		SELECT
		items.chrt_id,
		items.track_number,
		items.price,
		items.rid,
		items."name",
		items.sale,
		items.size,
		items.total_price,
		items.nm_id,
		items.brand,
		items.status
		FROM items 
		WHERE order_ref = $1
		`,
			uid,
		)
		if err != nil {
			return nil, err
		}
		var items []models.Item
		for rows.Next() {
			var i models.Item
			err = rows.Scan(
				&i.ChrtId,
				&i.TrackNumber,
				&i.Price,
				&i.Rid,
				&i.Name,
				&i.Sale,
				&i.Size,
				&i.TotalPrice,
				&i.NmId,
				&i.Brand,
				&i.Status,
			)

			if err != nil {
				return nil, err
			}
			items = append(items, i)
		}

		if rows.Err() != nil {
			return nil, rows.Err()
		}
	}
	return orders, rows.Err()
}

// Функция для добавления заказа в БД.
func (d *DB) Add(order models.Order) error {

	tx, err := d.db.Begin(context.Background())
	if err != nil {
		return err
	}
	_, err = tx.Exec(context.Background(), `
		INSERT INTO orders (
			order_uid,
			track_number,
			"entry",
			locale,
			internal_signature,
			customer_id,
			delivery_service,
			shardkey,
			sm_id,
			date_created,
			oof_shard)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);
	`,
		order.OrderUid,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSignature,
		order.CustomerId,
		order.DeliveryService,
		order.ShardKey,
		order.SmID,
		order.DateCreated,
		order.OffShard)
	if err != nil {
		tx.Rollback(context.Background())
		return err
	}

	_, err = tx.Exec(context.Background(), `
	INSERT INTO delivery (
		order_ref,
		"name",
		phone,
		zip,
		city,
		"address",
		region,
		email)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
	`,
		order.OrderUid,
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.Zip,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email)

	if err != nil {
		tx.Rollback(context.Background())
		return err
	}

	_, err = tx.Exec(context.Background(), `
		INSERT INTO payment (
			order_ref,
			"transaction",
			request_id,
			currency,
			"provider",
			amount,
			payment_dt,
			bank,
			delivery_cost,
			goods_total,
			custom_fee) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);
	`,
		order.OrderUid,
		order.Payment.Transaction,
		order.Payment.RequestId,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		order.Payment.PaymentDt,
		order.Payment.Bank,
		order.Payment.DeliveryCost,
		order.Payment.GoodsTotal,
		order.Payment.CustomFee)

	if err != nil {
		tx.Rollback(context.Background())
		return err
	}

	for _, o := range order.Items {
		_, err := tx.Exec(context.Background(), `
		INSERT INTO items (
		order_ref,
		chrt_id,
		track_number,
		price,
		rid,
		"name",
		sale,
		"size",
		total_price,
		nm_id,
		brand,
		"status") 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);
		`,
			order.OrderUid,
			o.ChrtId,
			o.TrackNumber,
			o.Price,
			o.Rid,
			o.Name,
			o.Sale,
			o.Size,
			o.TotalPrice,
			o.NmId,
			o.Brand,
			o.Status)
		if err != nil {
			tx.Rollback(context.Background())
			return err
		}
	}

	tx.Commit(context.Background())
	log.Printf("добавлен заказ с uid = %s", order.OrderUid)

	return nil
}
