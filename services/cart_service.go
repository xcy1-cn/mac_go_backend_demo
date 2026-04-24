package services

import (
	"demo/day6-9/db"
	"demo/day6-9/models"
)

func AddCart(cart models.Cart) (*models.Cart, error) {
	result, err := db.DB.Exec(
		"INSERT INTO carts (user_id, product_id, quantity) VALUES (?, ?, ?)",
		cart.UserID,
		cart.ProductID,
		cart.Quantity,
	)
	if err != nil {
		return nil, err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	cart.ID = int(lastID)
	return &cart, nil
}

func GetCartList(userID int64) ([]models.CartItemDetail, error) {
	rows, err := db.DB.Query(`
		SELECT 
			c.id,
			c.quantity,
			p.id,
			p.name,
			p.price,
			p.stock
		FROM carts c
		JOIN products p ON c.product_id = p.id
		WHERE c.user_id = ?
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var carts []models.CartItemDetail

	for rows.Next() {
		var item models.CartItemDetail
		err := rows.Scan(
			&item.ID,
			&item.Quantity,
			&item.Product.ID,
			&item.Product.Name,
			&item.Product.Price,
			&item.Product.Stock,
		)
		if err != nil {
			return nil, err
		}

		carts = append(carts, item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return carts, nil
}

func DeleteCart(id int, userID int64) (bool, error) {
	result, err := db.DB.Exec("DELETE FROM carts WHERE id = ? AND user_id = ?", id, userID)
	if err != nil {
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}

	if rowsAffected == 0 {
		return false, nil
	}

	return true, nil
}

func UpdateCartQuantity(id int, userID int64, quantity int) (*models.Cart, bool, error) {
	result, err := db.DB.Exec(
		"UPDATE carts SET quantity = ? WHERE id = ? AND user_id = ?",
		quantity,
		id,
		userID,
	)
	if err != nil {
		return nil, false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, false, err
	}

	if rowsAffected == 0 {
		return nil, false, nil
	}

	row := db.DB.QueryRow(
		"SELECT id, user_id, product_id, quantity FROM carts WHERE id = ? AND user_id = ?",
		id,
		userID,
	)

	var cart models.Cart
	err = row.Scan(&cart.ID, &cart.UserID, &cart.ProductID, &cart.Quantity)
	if err != nil {
		return nil, false, err
	}

	return &cart, true, nil
}

func GetCartSummary(userID int64) (*models.CartSummary, error) {
	row := db.DB.QueryRow(`
		SELECT 
			COUNT(*) as total_items,
			COALESCE(SUM(c.quantity),0) as total_quantity,
			COALESCE(SUM(c.quantity * p.price),0) as total_amount
		FROM carts c
		JOIN products p ON c.product_id = p.id
		WHERE c.user_id = ?
	`, userID)

	var summary models.CartSummary

	err := row.Scan(
		&summary.TotalItems,
		&summary.TotalQuantity,
		&summary.TotalAmount,
	)
	if err != nil {
		return nil, err
	}

	return &summary, nil
}
