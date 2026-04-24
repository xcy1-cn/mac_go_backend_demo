package services

import (
	"database/sql"
	"demo/day6-9/db"
	"demo/day6-9/models"
	"errors"
)

type cartOrderItem struct {
	ProductID int
	Quantity  int
	Price     float64
}

// CreateOrderFromCart
// 流程：
// 1. 查询当前用户购物车联表数据
// 2. 计算总金额
// 3. 插入当前用户 orders
// 4. 插入 order_items
// 5. 清空当前用户 carts
func CreateOrderFromCart(userID int64) (*models.OrderDetail, error) {
	tx, err := db.DB.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// 1. 查询当前用户购物车
	rows, err := tx.Query(`
		SELECT 
			c.product_id,
			c.quantity,
			p.price
		FROM carts c
		JOIN products p ON c.product_id = p.id
		WHERE c.user_id = ?
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cartItems []cartOrderItem
	var totalAmount float64

	for rows.Next() {
		var item cartOrderItem
		err = rows.Scan(&item.ProductID, &item.Quantity, &item.Price)
		if err != nil {
			return nil, err
		}

		cartItems = append(cartItems, item)
		totalAmount += float64(item.Quantity) * item.Price
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(cartItems) == 0 {
		return nil, errors.New("cart is empty")
	}

	// 2. 插入 orders 主表（绑定 user_id）
	orderResult, err := tx.Exec(
		"INSERT INTO orders (user_id, total_amount, status) VALUES (?, ?, ?)",
		userID,
		totalAmount,
		"created",
	)
	if err != nil {
		return nil, err
	}

	orderID64, err := orderResult.LastInsertId()
	if err != nil {
		return nil, err
	}
	orderID := int(orderID64)

	// 3. 插入 order_items 明细表
	var orderItems []models.OrderItem

	for _, item := range cartItems {
		itemResult, execErr := tx.Exec(
			"INSERT INTO order_items (order_id, product_id, quantity, price) VALUES (?, ?, ?, ?)",
			orderID,
			item.ProductID,
			item.Quantity,
			item.Price,
		)
		if execErr != nil {
			err = execErr
			return nil, err
		}

		itemID64, lastErr := itemResult.LastInsertId()
		if lastErr != nil {
			err = lastErr
			return nil, err
		}

		orderItems = append(orderItems, models.OrderItem{
			ID:        int(itemID64),
			OrderID:   orderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		})
	}

	// 4. 清空当前用户购物车
	_, err = tx.Exec("DELETE FROM carts WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}

	// 5. 提交事务
	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	// 6. 查询订单主表（按 id + user_id）
	var order models.Order
	row := db.DB.QueryRow(`
		SELECT id, user_id, total_amount, status, created_at
		FROM orders
		WHERE id = ? AND user_id = ?
	`, orderID, userID)

	err = row.Scan(
		&order.ID,
		&order.UserID,
		&order.TotalAmount,
		&order.Status,
		&order.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("order created but not found")
		}
		return nil, err
	}

	detail := &models.OrderDetail{
		ID:          order.ID,
		UserID:      order.UserID,
		TotalAmount: order.TotalAmount,
		Status:      order.Status,
		CreatedAt:   order.CreatedAt,
		Items:       orderItems,
	}

	return detail, nil
}

// GetOrders
// 查询当前用户订单列表
func GetOrders(userID int64) ([]models.Order, error) {
	rows, err := db.DB.Query(`
		SELECT id, user_id, total_amount, status, created_at
		FROM orders
		WHERE user_id = ?
		ORDER BY id DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order

	for rows.Next() {
		var order models.Order
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.TotalAmount,
			&order.Status,
			&order.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

// GetOrderByID
// 查询当前用户单个订单详情：主表 + 明细表
func GetOrderByID(id int, userID int64) (*models.OrderDetail, error) {
	// 1. 查订单主表，必须带 user_id
	row := db.DB.QueryRow(`
		SELECT id, user_id, total_amount, status, created_at
		FROM orders
		WHERE id = ? AND user_id = ?
	`, id, userID)

	var order models.OrderDetail
	err := row.Scan(
		&order.ID,
		&order.UserID,
		&order.TotalAmount,
		&order.Status,
		&order.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// 2. 查订单明细
	rows, err := db.DB.Query(`
		SELECT id, order_id, product_id, quantity, price
		FROM order_items
		WHERE order_id = ?
		ORDER BY id ASC
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.OrderItem

	for rows.Next() {
		var item models.OrderItem
		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.Quantity,
			&item.Price,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	order.Items = items
	return &order, nil
}
