package controller

import (
	"database/shopify"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiResponsePriceRules struct {
	PriceRules []shopify.PriceRule `json:"pricerules"`
	Message    string              `json:"message"`
}

func PriceRules(c *gin.Context) {
	price_rules, err := shopify.FetchPriceRules()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ApiResponsePriceRules{
			PriceRules: nil,
			Message:    fmt.Sprintf("Error: %v", err),
		})
		return
	}
	c.JSON(http.StatusOK, ApiResponsePriceRules{
		PriceRules: price_rules,
		Message:    "Price Rules fetched successfully",
	})
}
