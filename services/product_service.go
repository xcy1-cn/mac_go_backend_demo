package services

import (
	"database/sql"
	"demo/day6-9/db"
	"demo/day6-9/models"
)

// 获取全部
func GetAllProducts() ([]models.Product, error) {
	rows, err := db.DB.Query("SELECT id, name, price, stock FROM products")
	if err = rows.Err(); err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product

	for rows.Next() {
		var product models.Product

		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.Stock,
		)

		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, nil
}

// 根据ID查找
func GetProductByID(id int) (*models.Product, error) {
	row := db.DB.QueryRow("SELECT id, name, price, stock FROM products WHERE id = ?", id)

	var product models.Product
	err := row.Scan(&product.ID, &product.Name, &product.Price, &product.Stock)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &product, nil
}

// 添加
func AddProduct(product models.Product) (*models.Product, error) {
	result, err := db.DB.Exec(
		"INSERT INTO products (name, price, stock) VALUES (?, ?, ?)",
		product.Name,
		product.Price,
		product.Stock,
	)
	if err != nil {
		return nil, err
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	product.ID = int(lastID)
	return &product, nil
}

// 删除
func DeleteProduct(id int) (bool, error) {
	result, err := db.DB.Exec("DELETE FROM products WHERE id = ?", id)
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

// 更新
func UpdateProduct(id int, newProduct models.Product) (*models.Product, bool, error) {
	result, err := db.DB.Exec(
		"UPDATE products SET name = ?, price = ?, stock = ? WHERE id = ?",
		newProduct.Name,
		newProduct.Price,
		newProduct.Stock,
		id,
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

	updatedProduct := models.Product{
		ID:    id,
		Name:  newProduct.Name,
		Price: newProduct.Price,
		Stock: newProduct.Stock,
	}

	return &updatedProduct, true, nil
}
