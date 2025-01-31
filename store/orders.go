package store

import (
	"context"
	"database/sql"

	"github.com/lib/pq"
	"github.com/poohda-go/types"
)

type OrdersStore struct {
	db *sql.DB
}

func (s *OrdersStore) GetAllOrders() ([]types.Order, error) {
	orders := []types.Order{}
	query := `SELECT o.id, o.name, o.quantity, o.address, o.price, o.is_delivered, array_agg(cb.clothe_id) AS "clothes_ordered" FROM "orders" AS o JOIN "clothes_bought" AS "cb" ON cb.order_id = o.id GROUP BY o.id, o.name, o.quantity, o.address, o.price, o.is_delivered`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var order types.Order

		if err := rows.Scan(
			&order.Id,
			&order.Name,
			&order.Quantity,
			&order.Address,
			&order.Price,
			&order.IsDelivered,
			pq.Array(&order.ClothesBought),
		); err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}

func (s *OrdersStore) GetASingleOrder(ctx context.Context, id int) (*types.Order, error) {
	var order types.Order
	query := `SELECT o.id, o.name, o.quantity, o.address, o.price, o.is_delivered, array_agg(cb.clothe_id) AS "clothes_ordered" FROM "orders" AS o JOIN "clothes_bought" AS "cb" ON cb.order_id = o.id WHERE id=$1 GROUP BY o.id, o.name, o.quantity, o.address, o.price, o.is_delivered`

	err := s.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&order.Id,
		&order.Name,
		&order.Quantity,
		&order.Address,
		&order.Price,
		&order.IsDelivered,
		pq.Array(&order.ClothesBought),
	)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (s *OrdersStore) CreateANewOrder(ctx context.Context, payload types.OrderDTO) (*types.Order, error) {
	var order types.Order
	var clothesBought int
	query := `INSERT INTO "orders" (name, quantity, address, price, is_delivered) VALUES ($1, $2, $3, $4, $5) RETURNING name, quantity, address, price, is_delivered`
	clothesBoughtQuery := `INSERT INTO "clothes_bought" (order_id, clothe_id, quantity) VALUES ($1, $2, $3) RETURNING id`

	err := s.db.QueryRowContext(
		ctx,
		query,
		payload.Name,
		payload.Quantity,
		payload.Address,
		payload.Price,
		payload.IsDelivered,
	).Scan(
		&order.Name,
		&order.Quantity,
		&order.Address,
		&order.Price,
		&order.IsDelivered,
	)
	if err != nil {
		return nil, err
	}

	for _, clotheBought := range payload.ClothesBought {
		err := s.db.QueryRowContext(ctx, clothesBoughtQuery, &order.Id, clotheBought.Id, clotheBought.Quantity).Scan(
			&clothesBought,
		)
		if err != nil {
			return nil, err
		}
	}

	return &order, nil
}
