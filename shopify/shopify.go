package shopify

import (
	"bytes"
	"database/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Product struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Currency string `json:"currency"`
	Vendor   string `json:"vendor"`
	Images   []struct {
		ID  int    `json:"id"`
		Src string `json:"src"`
	} `json:"images"`
	Variants []struct {
		Id         int    `json:"id"`
		ProductId  int    `json:"product_id"`
		Title      string `json:"title"`
		Price      string `json:"price"`
		ImageId    int    `json:"image_id"`
		VariantIDs []int  `json:"variant_ids"`
	} `json:"variants"`
}

type CheckOut struct {
	ID        int                    `json:"id"`
	CartToken string                 `json:"cart_token"`
	Customer  map[string]interface{} `json:"customer"`
	LineItems []interface{}          `json:"line_items"`
}

type PriceRule struct {
	ID       int    `json:"id"`
	StartsAt string `json:"starts_at"`
	EndsAt   string `json:"ends_at"`
	Value    string `json:"value"`
}

type DiscountCode struct {
	ID          int    `json:"id"`
	PriceRuleId int    `json:"price_rule_id"`
	Code        string `json:"code"`
	UsageCount  int    `json:"usage_count"`
}

type Customer struct {
	ID    int    `json:"id"`
	State string `json:"state"`
}

type Order struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Currency string `json:"currency"`
	Customer struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
	} `json:"customer"`
}

type ShopifyResponse struct {
	Products      []Product      `json:"products"`
	Orders        []Order        `json:"orders"`
	Customers     []Customer     `json:"customers"`
	DraftOrders   []DraftOrder   `json:"draft_orders"`
	CheckOuts     []CheckOut     `json:"checkouts"`
	PriceRules    []PriceRule    `json:"price_rules"`
	DiscountCodes []DiscountCode `json:"discount_codes"`
}

type DraftOrder struct {
	LineItems []LineItem `json:"line_items"`
	Customer  *Customer  `json:"customer,omitempty"`
}

type LineItem struct {
	VariantID int    `json:"variant_id"`
	Quantity  int    `json:"quantity"`
	Title     string `json:"title"`
}

func CreateDraftOrder(cart DraftOrder) (*DraftOrder, error) {
	url := fmt.Sprintf("https://%s:%s@%s.myshopify.com/admin/api/%s/draft_orders.json",
		config.ShopifyStore.APIKey,
		config.ShopifyStore.Password,
		config.ShopifyStore.ShopName,
		config.ShopifyStore.APIVersion)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	data, err := json.Marshal(map[string]DraftOrder{"draft_order": cart})
	if err != nil {
		return nil, fmt.Errorf("error marshaling draft order: %v", err)
	}

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("error making request to Shopify API: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != 201 {
		return nil, fmt.Errorf("non-201 status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	draftOrder := result["draft_order"].(map[string]interface{})
	var order DraftOrder
	orderJSON, _ := json.Marshal(draftOrder)
	err = json.Unmarshal(orderJSON, &order)

	if err != nil {
		return nil, fmt.Errorf("error unmarshaling order: %v", err)
	}

	return &order, nil
}

func FetchDraftOrder() ([]DraftOrder, error) {
	url := fmt.Sprintf("https://%s:%s@%s.myshopify.com/admin/api/%s/draft_orders.json",
		config.ShopifyStore.APIKey,
		config.ShopifyStore.Password,
		config.ShopifyStore.ShopName,
		config.ShopifyStore.APIVersion)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making request to Shopify API: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("non-200 status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var shopifyResp ShopifyResponse
	err = json.Unmarshal(body, &shopifyResp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	return shopifyResp.DraftOrders, nil
}

func FetchProducts() ([]Product, error) {
	url := fmt.Sprintf("https://%s:%s@%s.myshopify.com/admin/api/%s/products.json",
		config.ShopifyStore.APIKey,
		config.ShopifyStore.Password,
		config.ShopifyStore.ShopName,
		config.ShopifyStore.APIVersion)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making request to Shopify API: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("non-200 status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var shopifyResp ShopifyResponse
	err = json.Unmarshal(body, &shopifyResp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	return shopifyResp.Products, nil
}

func FetchCheckouts() ([]CheckOut, error) {
	url := fmt.Sprintf("https://%s:%s@%s.myshopify.com/admin/api/%s/checkouts.json",
		config.ShopifyStore.APIKey,
		config.ShopifyStore.Password,
		config.ShopifyStore.ShopName,
		config.ShopifyStore.APIVersion)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making request to Shopify API: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("non-200 status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var shopifyResp ShopifyResponse
	err = json.Unmarshal(body, &shopifyResp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	return shopifyResp.CheckOuts, nil
}

func FetchPriceRules() ([]PriceRule, error) {
	url := fmt.Sprintf("https://%s:%s@%s.myshopify.com/admin/api/%s/price_rules.json",
		config.ShopifyStore.APIKey,
		config.ShopifyStore.Password,
		config.ShopifyStore.ShopName,
		config.ShopifyStore.APIVersion)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making request to Shopify API: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("non-200 status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var shopifyResp ShopifyResponse
	err = json.Unmarshal(body, &shopifyResp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	return shopifyResp.PriceRules, nil
}

func FetchDiscountcodes(d int64) ([]DiscountCode, error) {
	url := fmt.Sprintf("https://%s:%s@%s.myshopify.com/admin/api/%s/price_rules/%d/discount_codes.json",
		config.ShopifyStore.APIKey,
		config.ShopifyStore.Password,
		config.ShopifyStore.ShopName,
		config.ShopifyStore.APIVersion, d)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making request to Shopify API: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("non-200 status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var shopifyResp ShopifyResponse
	err = json.Unmarshal(body, &shopifyResp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	return shopifyResp.DiscountCodes, nil
}

func FetchCustomers() ([]Customer, error) {
	url := fmt.Sprintf("https://%s:%s@%s.myshopify.com/admin/api/%s/customers.json",
		config.ShopifyStore.APIKey,
		config.ShopifyStore.Password,
		config.ShopifyStore.ShopName,
		config.ShopifyStore.APIVersion)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making request to Shopify API: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("non-200 status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var shopifyResp ShopifyResponse
	err = json.Unmarshal(body, &shopifyResp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	return shopifyResp.Customers, nil
}

func FetchOrders() ([]Order, error) {
	url := fmt.Sprintf("https://%s:%s@%s.myshopify.com/admin/api/%s/orders.json",
		config.ShopifyStore.APIKey,
		config.ShopifyStore.Password,
		config.ShopifyStore.ShopName,
		config.ShopifyStore.APIVersion)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making request to Shopify API: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("non-200 status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var shopifyResp ShopifyResponse
	err = json.Unmarshal(body, &shopifyResp)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	return shopifyResp.Orders, nil
}

func RemoveDraftOrder(draftOrderID int) error {

	url := fmt.Sprintf("https://%s:%s@%s.myshopify.com/admin/api/%s/draft_orders/%d.json",
		config.ShopifyStore.APIKey,
		config.ShopifyStore.Password,
		config.ShopifyStore.ShopName,
		config.ShopifyStore.APIVersion,
		draftOrderID,
	)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request to Shopify API: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("non-200 status code: %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}
