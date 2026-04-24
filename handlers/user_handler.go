package handlers

import (
	"demo/day6-9/models"
	"demo/day6-9/services"
	"demo/day6-9/utils"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var req models.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "invalid request body")
		return
	}

	err := services.RegisterUser(req)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	utils.Success(c, "register success")
}

func Login(c *gin.Context) {
	var req models.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, 400, "invalid request body")
		return
	}

	user, err := services.LoginUser(req)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		utils.Error(c, 500, "failed to generate token")
		return
	}

	utils.Success(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"nickname": user.Nickname,
		"token":    token,
	})
}

func GetUserInfo(c *gin.Context) {
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

	user, err := services.GetUserByID(userID)
	if err != nil {
		utils.Error(c, 500, "failed to get user info")
		return
	}

	utils.Success(c, gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"nickname":   user.Nickname,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	})
}
