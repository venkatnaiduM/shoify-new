package controller

import (
	"bytes"
	"database/config"
	"database/shopify"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AddToCartRequest struct {
	VariantID int    `form:"variant_id" binding:"required"`
	Quantity  int    `form:"quantity" binding:"required"`
	Title     string `form:"title" binding:"required"`
}

type ApiResponseDraftOrders struct {
	DraftOrder []shopify.DraftOrder `json:"draft_orders"`
	Message    string               `json:"message"`
}

type ApiResponseCart struct {
	ShopifyRes ShopifyResponse `json:"shopifyResponse"`
	Message    string          `json:"message"`
}
type ShopifyResponse struct {
	Data struct {
		Cart struct {
			ID    string `json:"id"`
			Lines struct {
				Edges []struct {
					Node struct {
						ID          string `json:"id"`
						Quantity    int    `json:"quantity"`
						Merchandise struct {
							ID    string `json:"id"`
							Title string `json:"title"`
						} `json:"merchandise"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"lines"`
		} `json:"cart"`
	} `json:"data"`
}

func GetCartDetails(c *gin.Context) {
	cartID := c.DefaultQuery("cart_id", "")
	if cartID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart ID is required"})
		return
	}

	query := `
			query cart($id: ID!) {
				cart(id: $id) {
					id
					lines(first: 10) {
						edges {
							node {
								id
								quantity
								merchandise {
									... on ProductVariant {
										id
										title
									}
								}
							}
						}
					}
				}
			}
		`

	requestBody := map[string]interface{}{
		"query": query,
		"variables": map[string]interface{}{
			"id": cartID,
		},
	}

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", config.ShopifyStore.Url, bytes.NewBuffer(requestBodyJSON))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Shopify-Storefront-Access-Token", config.ShopifyStore.StorefrontAccessToken)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		return
	}

	var shopifyResponse ShopifyResponse
	err = json.Unmarshal(body, &shopifyResponse)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, ApiResponseCart{
		ShopifyRes: shopifyResponse,
		Message:    "Cart Details Fetched successfully",
	})
}

func DeleteCartDetails(c *gin.Context) {
	cartID := c.DefaultQuery("cart_id", "")
	if cartID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart ID is required"})
		return
	}

	lineID := c.DefaultQuery("line_id", "")
	if lineID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Line ID is required"})
		return
	}

	query := `mutation cartLinesRemove($cartId: ID!, $lineIds: [ID!]!) { 
	cartLinesRemove(cartId: $cartId, lineIds: $lineIds) {
	 cart {
	  id lines(first: 10) { edges { node { id quantity } } } }} }`
	requestBody := map[string]interface{}{
		"query": query,
		"variables": map[string]interface{}{
			"cartId":  cartID,
			"lineIds": []string{lineID},
		},
	}

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", config.ShopifyStore.Url, bytes.NewBuffer(requestBodyJSON))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Shopify-Storefront-Access-Token", config.ShopifyStore.StorefrontAccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		return
	}
	var shopifyResponse ShopifyResponse
	err = json.Unmarshal(body, &shopifyResponse)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, ApiResponseCart{
		ShopifyRes: shopifyResponse,
		Message:    "Cart deleted successfully",
	})

}

func CartLinesAdd(c *gin.Context) {
	cartID := c.DefaultQuery("cart_id", "")
	if cartID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart ID is required"})
		return
	}

	merchandiseID := c.DefaultQuery("merchent_id", "")
	if merchandiseID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Line ID is required"})
		return
	}

	query := `mutation cartLinesAdd($id: ID!, $lines: [CartLineInput!]!) {
	 cartLinesAdd(cartId: $id, lines: $lines) { 
	 cart { 
	 id lines(first: 5) {
	  edges { node { id quantity merchandise { 
	  ... on ProductVariant { id title priceV2
	  { amount currencyCode } } } } } } } } }`
	requestBody := map[string]interface{}{
		"query": query,
		"variables": map[string]interface{}{
			"id": cartID,
			"lines": []interface{}{
				map[string]interface{}{
					"quantity":      2,
					"merchandiseId": merchandiseID,
				},
			},
		},
	}

	requestBodyJSON, err := json.Marshal(requestBody)
	if err != nil {
		return
	}
	req, err := http.NewRequest("POST", config.ShopifyStore.Url, bytes.NewBuffer(requestBodyJSON))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Shopify-Storefront-Access-Token", config.ShopifyStore.StorefrontAccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		return
	}
	var shopifyResponse ShopifyResponse
	err = json.Unmarshal(body, &shopifyResponse)
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, ApiResponseCart{
		ShopifyRes: shopifyResponse,
		Message:    "Cart Details Fetched successfully",
	})

}

func AddToCart(c *gin.Context) {
	var request AddToCartRequest
	if err := c.ShouldBind(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data", "message": err.Error()})
		return
	}
	lineItem := shopify.LineItem{
		VariantID: request.VariantID,
		Quantity:  request.Quantity,
		Title:     request.Title,
	}

	draftOrder := shopify.DraftOrder{
		LineItems: []shopify.LineItem{lineItem},
	}
	createdOrder, err := shopify.CreateDraftOrder(draftOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create draft order: %v", err)})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Product added to cart",
		"order":   createdOrder,
	})
}

func GetCart(c *gin.Context) {
	draftorders, err := shopify.FetchDraftOrder()
	if err != nil {
		c.JSON(http.StatusInternalServerError, ApiResponseDraftOrders{
			DraftOrder: nil,
			Message:    fmt.Sprintf("Error: %v", err),
		})
		return
	}
	c.JSON(http.StatusOK, ApiResponseDraftOrders{
		DraftOrder: draftorders,
		Message:    "Orders fetched successfully",
	})
}

func RemoveFromCart(c *gin.Context) {

	variantIDStr := c.PostForm("variant_id")
	if variantIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Variant ID is required"})
		return
	}

	variantID := 0
	_, err := fmt.Sscanf(variantIDStr, "%d", &variantID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Variant ID"})
		return
	}

	err = shopify.RemoveDraftOrder(variantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to remove item: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product removed from cart successfully"})
}
