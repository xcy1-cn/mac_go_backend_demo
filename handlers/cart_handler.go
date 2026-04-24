package handlers

import (
	"demo/day6-9/models"
	"demo/day6-9/services"
	"demo/day6-9/utils"
	"fmt"

	"github.com/gin-gonic/gin"
)

func AddCart(c *gin.Context) {
	userIdValue, exists := c.Get("userId")
	fmt.Println("handler userIdValue:", userIdValue, "exists:", exists)

	var cart models.Cart
	err := c.ShouldBindJSON(&cart)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	// userIdValue, exists := c.Get("userId")
	fmt.Println("handler userIdValue:", userIdValue, "exists:", exists)
	if !exists {
		utils.Error(c, 401, "user not found in context")
		return
	}

	userID, ok := userIdValue.(int64)
	if !ok {
		utils.Error(c, 500, "invalid user id type")
		return
	}

	cart.UserID = userID

	newCart, err := services.AddCart(cart)
	if err != nil {
		fmt.Println("AddCart error:", err)
		utils.Error(c, 500, err.Error())
		return
	}

	utils.Success(c, newCart)
}

func GetCartList(c *gin.Context) {
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

	carts, err := services.GetCartList(userID)
	if err != nil {
		utils.Error(c, 500, "query cart failed")
		return
	}

	utils.Success(c, carts)
}

func DeleteCartByID(c *gin.Context) {
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

	okDeleted, err := services.DeleteCart(id, userID)
	if err != nil {
		utils.Error(c, 500, "delete cart failed")
		return
	}

	if !okDeleted {
		utils.Error(c, 404, "cart not found")
		return
	}

	utils.Success(c, gin.H{"message": "deleted"})
}

func UpdateCartByID(c *gin.Context) {
	id, err := GetIDFromParam(c)
	if err != nil {
		utils.Error(c, 400, "invalid id")
		return
	}

	var req models.UpdateCartReq
	err = c.ShouldBindJSON(&req)
	if err != nil {
		utils.Error(c, 400, err.Error())
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

	updated, found, err := services.UpdateCartQuantity(id, userID, req.Quantity)
	if err != nil {
		utils.Error(c, 500, "update cart failed")
		return
	}

	if !found {
		utils.Error(c, 404, "cart not found")
		return
	}

	utils.Success(c, updated)
}

func GetCartSummary(c *gin.Context) {
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

	data, err := services.GetCartSummary(userID)
	if err != nil {
		utils.Error(c, 500, "summary failed")
		return
	}

	utils.Success(c, data)
}
