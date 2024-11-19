package controller

import (
	"database/shopify"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiResponseProducts struct {
	Products []shopify.Product `json:"products"`
	Message  string            `json:"message"`
}

func ProductDetails(c *gin.Context) {
	products, err := shopify.FetchProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ApiResponseProducts{
			Products: nil,
			Message:  fmt.Sprintf("Error: %v", err),
		})
		return
	}
	c.JSON(http.StatusOK, ApiResponseProducts{
		Products: products,
		Message:  "Products fetched successfully",
	})
}
