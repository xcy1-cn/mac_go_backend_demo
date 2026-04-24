package handlers

import (
	"demo/day6-9/models"
	"demo/day6-9/services"
	"demo/day6-9/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 获取全部商品
func GetProducts(c *gin.Context) {
	products, err := services.GetAllProducts()
	if err != nil {
		utils.Error(c, 500, "query failed")
		return
	}

	utils.Success(c, products)
}

// 获取单个商品--id
func GetProductByID(c *gin.Context) {
	id, err := GetIDFromParam(c)
	if err != nil {
		// c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		utils.Error(c, 400, "invalid id")
		return
	}

	product, err := services.GetProductByID(id)
	if err != nil {
		utils.Error(c, 500, "query failed")
		return
	}

	if product == nil {
		utils.Error(c, 404, "product not found")
		return
	}

	utils.Success(c, product)
}

// 添加一个商品
func AddProduct(c *gin.Context) {
	var product models.Product
	err := c.ShouldBindJSON(&product)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	newProduct, err := services.AddProduct(product)
	if err != nil {
		utils.Error(c, 500, "insert failed")
		return
	}

	utils.Success(c, newProduct)
}

// 删除一个商品--id
func DeleteProductByID(c *gin.Context) {
	id, err := GetIDFromParam(c)
	if err != nil {
		utils.Error(c, 400, "invalid id")
		return
	}

	ok, err := services.DeleteProduct(id)
	if err != nil {
		utils.Error(c, 500, "delete failed")
		return
	}

	if !ok {
		utils.Error(c, 404, "product not found")
		return
	}

	utils.Success(c, gin.H{"message": "deleted"})
}

// 更新一个商品--id
func UpdateProductByID(c *gin.Context) {
	id, err := GetIDFromParam(c)
	if err != nil {
		utils.Error(c, 400, "invalid id")
		return
	}

	var product models.Product
	err = c.ShouldBindJSON(&product)
	if err != nil {
		utils.Error(c, 400, err.Error())
		return
	}

	updated, ok, err := services.UpdateProduct(id, product)
	if err != nil {
		utils.Error(c, 500, "update failed")
		return
	}

	if !ok {
		utils.Error(c, 404, "product not found")
		return
	}

	utils.Success(c, updated)
}

func GetIDFromQuery(c *gin.Context) (int, error) {
	idStr := c.Query("id")
	return strconv.Atoi(idStr)
}

func GetIDFromParam(c *gin.Context) (int, error) {
	idStr := c.Param("id")
	return strconv.Atoi(idStr)
}
