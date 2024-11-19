package controller

import (
	"database/shopify"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiResponseOrders struct {
	Orders  []shopify.Order `json:"orders"`
	Message string          `json:"message"`
}

func OrderDetails(c *gin.Context) {
	orders, err := shopify.FetchOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ApiResponseOrders{
			Orders:  nil,
			Message: fmt.Sprintf("Error: %v", err),
		})
		return
	}
	c.JSON(http.StatusOK, ApiResponseOrders{
		Orders:  orders,
		Message: "Orders fetched successfully",
	})
}
