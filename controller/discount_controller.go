package controller

import (
	"database/shopify"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ApiResponseDiscount struct {
	DiscountCodes []shopify.DiscountCode `json:"discount_codes"`
	Message       string                 `json:"message"`
}

func DiscountCodes(c *gin.Context) {

	discountCode := c.DefaultQuery("price_rule_id", "")
	if discountCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart ID is required"})
		return
	}

	DiscountCode, err := strconv.ParseInt(discountCode, 10, 64)
	if err != nil {
		log.Println("Error parsing DiscountCode:", err)
		c.JSON(400, gin.H{"error": "Invalid Discount"})
		return
	}
	// var DiscountCode int64 = 1452060672282

	discount_codes, err := shopify.FetchDiscountcodes(DiscountCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ApiResponseDiscount{
			DiscountCodes: nil,
			Message:       fmt.Sprintf("Error: %v", err),
		})
		return
	}
	c.JSON(http.StatusOK, ApiResponseDiscount{
		DiscountCodes: discount_codes,
		Message:       "Discount Codes fetched successfully",
	})
}
