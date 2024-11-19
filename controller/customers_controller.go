package controller

import (
	"database/shopify"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiResponseCustomers struct {
	Customers []shopify.Customer `json:"customers"`
	Message   string             `json:"message"`
}

func CustomerDetails(c *gin.Context) {
	customers, err := shopify.FetchCustomers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ApiResponseCustomers{
			Customers: nil,
			Message:   fmt.Sprintf("Error: %v", err),
		})
		return
	}
	c.JSON(http.StatusOK, ApiResponseCustomers{
		Customers: customers,
		Message:   "Orders fetched successfully",
	})
}
