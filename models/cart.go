package models

type Cart struct {
	ID int `json:"id"`
	// userId
	UserID    int64 `json:"user_id"`
	ProductID int   `json:"product_id" binding:"required"`
	Quantity  int   `json:"quantity" binding:"gte=1"`
}

type UpdateCartReq struct {
	Quantity int `json:"quantity" binding:"gte=1"`
}

type CartItemDetail struct {
	ID       int     `json:"id"`
	Quantity int     `json:"quantity"`
	Product  Product `json:"product"`
}

type CartSummary struct {
	TotalItems    int     `json:"total_items"`
	TotalQuantity int     `json:"total_quantity"`
	TotalAmount   float64 `json:"total_amount"`
}
