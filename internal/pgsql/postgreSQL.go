package pgsql

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Artamonov-Georgii/l0/internal/server"
	_ "github.com/lib/pq"
)

var Db *sql.DB
var err_db error

func StartSql() {
	Db, err_db = sql.Open("postgres", "postgres://postgres:postgrespw@localhost:32768/postgres?sslmode=disable")

	if err_db != nil {
		log.Fatal("Not possible to connect to the database")
	}
}

func InsertOrder(db *sql.DB, o server.Order) error {

	deliveryJson, err := json.Marshal(o.Delivery)
	if err != nil {
		return fmt.Errorf("failed to marshal delivery data: %v", err)
	}

	paymentJson, err := json.Marshal(o.Payment)
	if err != nil {
		return fmt.Errorf("failed to marshal payment data: %v", err)
	}

	query := `
        INSERT INTO orders (
            order_uid, track_number, entry, delivery, payment,
            locale, internal_signature, customer_id, delivery_service,
            shardkey, sm_id, date_created, oof_shard
        )
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
    `

	_, err = db.Exec(query, o.OrderUID, o.TrackNum, o.Entry, string(deliveryJson), string(paymentJson),
		o.Locale, "", o.Customer, o.Service, o.ShardKey, o.SMID, o.CreatedAt, o.OOFShard)

	if err != nil {
		return fmt.Errorf("failed to insert order data: %v", err)
	}

	for _, item := range o.Items {
		query := `
            INSERT INTO order_items (
                order_id, chrt_id, track_number, price, rid, name,
                sale, size, total_price, nm_id, brand, status
            )
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
        `
		_, err = db.Exec(query, o.OrderUID, item.ChrtID, item.TrackNum, item.Price, item.RID,
			item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			return fmt.Errorf("failed to insert order item data: %v", err)
		}
	}

	return nil
}

func GetAllOrders(db *sql.DB) (map[string]server.Order, error) {
    orders := make(map[string]server.Order)

    rows, err := db.Query("SELECT * FROM orders")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var order server.Order
        var deliveryJson, paymentJson []byte

        err := rows.Scan(&order.OrderUID, &order.TrackNum, &order.Entry, &deliveryJson, &paymentJson, &order.Locale, &order.IntSig,
            &order.Customer, &order.Service, &order.ShardKey, &order.SMID, &order.CreatedAt, &order.OOFShard)

        if err != nil {
            return nil, err
        }

        if err := json.Unmarshal(deliveryJson, &order.Delivery); err != nil {
            return nil, err
        }

        if err := json.Unmarshal(paymentJson, &order.Payment); err != nil {
            return nil, err
        }

        orderItems := []server.Item{}
        itemRows, err := db.Query(`SELECT oi.chrt_id, oi.track_number, oi.price, 
		oi.rid, oi.name, oi.sale, 
		oi.size, oi.total_price, oi.nm_id, 
		oi.brand, oi.status
		FROM order_items oi
		WHERE oi.order_id = $1`, order.OrderUID)
        if err != nil {
            return nil, err
        }
        defer itemRows.Close()

        for itemRows.Next() {
            var item server.Item
            err := itemRows.Scan(&item.ChrtID, &item.TrackNum, &item.Price, 
				&item.RID, &item.Name, &item.Sale, 
				&item.Size, &item.TotalPrice, &item.NmID, 
				&item.Brand, &item.Status)

            if err != nil {
                return nil, err
            }
            orderItems = append(orderItems, item)
        }

        if err := itemRows.Err(); err != nil {
            return nil, err
        }

        order.Items = orderItems

        orders[order.OrderUID] = order
    }

    if err := rows.Err(); err != nil {
        return nil, err
    }

    return orders, nil
}




