package handlers

import (
	"demo/day6-9/services"
	"demo/day6-9/utils"

	"github.com/gin-gonic/gin"
)

func CreateOrder(c *gin.Context) {
	userIdValue, exists := c.Get("userId")
	if !exists {
		utils.Error(c, 401, "user not found in context")
		return
	}

	userID, ok := userIdValue.(int64)
	if !ok {
		utils.Error(c, 500, "invalid user id type")
		return
	}

	order, err := services.CreateOrderFromCart(userID)
	if err != nil {
		if err.Error() == "cart is empty" {
			utils.Error(c, 400, "cart is empty")
			return
		}

		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, order)
}

func GetOrders(c *gin.Context) {
	userIdValue, exists := c.Get("userId")
	if !exists {
		utils.Error(c, 401, "user not found in context")
		return
	}

	userID, ok := userIdValue.(int64)
	if !ok {
		utils.Error(c, 500, "invalid user id type")
		return
	}

	orders, err := services.GetOrders(userID)
	if err != nil {
		utils.Error(c, 500, "query orders failed")
		return
	}

	utils.Success(c, orders)
}

func GetOrderByID(c *gin.Context) {
	id, err := GetIDFromParam(c)
	if err != nil {
		utils.Error(c, 400, "invalid id")
		return
	}

	userIdValue, exists := c.Get("userId")
	if !exists {
		utils.Error(c, 401, "user not found in context")
		return
	}

	userID, ok := userIdValue.(int64)
	if !ok {
		utils.Error(c, 500, "invalid user id type")
		return
	}

	order, err := services.GetOrderByID(id, userID)
	if err != nil {
		utils.Error(c, 500, "query order failed")
		return
	}

	if order == nil {
		utils.Error(c, 404, "order not found")
		return
	}

	utils.Success(c, order)
}
