package controller

import (
	"database/shopify"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiResponseCheckouts struct {
	CheckOuts []shopify.CheckOut `json:"checkouts"`
	Message   string             `json:"message"`
}

func CheckOutDetails(c *gin.Context) {
	checkouts, err := shopify.FetchCheckouts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ApiResponseCheckouts{
			CheckOuts: nil,
			Message:   fmt.Sprintf("Error: %v", err),
		})
		return
	}
	c.JSON(http.StatusOK, ApiResponseCheckouts{
		CheckOuts: checkouts,
		Message:   "Products fetched successfully",
	})
}
