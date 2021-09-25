package controllers

import (
	"net/http"
	"restapi/models"
	"restapi/schemas"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func AddStock(c *gin.Context) {
	var payload schemas.AddStockSchema

	// validation
	if err := c.ShouldBindWith(&payload, binding.JSON); err != nil {
		_ = c.AbortWithError(422, err).SetType(gin.ErrorTypeBind)
		return
	}

	// check drug is exists
	drug := models.CheckDrugName(payload.Name)
	if drug != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "drug not found."})
		return
	}

	// check price must less than max price
	max_price := models.GetMaxPrice(payload.Name)
	if payload.Price > max_price {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "the price must be below the max price"})
		return
	}

	// add stock
	stock := models.GetStockLogistics(payload.Name)
	payload.Stock = payload.Stock + stock

	models.CreateNewStock(payload)
	c.JSON(http.StatusOK, gin.H{"detail": "Successfully add new stock."})
}
