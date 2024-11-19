package controller

import (
	"bytes"
	"database/config"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"strconv"

	"github.com/gin-gonic/gin"
)

func HandleOrderSubmission(c *gin.Context) {

	variantId := c.PostForm("variant_id")
	quantity := c.PostForm("quantity")
	customerId := c.PostForm("customer[id]")
	CustomerEmail := c.PostForm("customer[email]")
	Address1 := c.PostForm("shipping_address[address1]")
	City := c.PostForm("shipping_address[city]")
	Province := c.PostForm("shipping_address[province]")
	Country := c.PostForm("shipping_address[country]")
	FinancialStatus := c.PostForm("financial_status")
	fmt.Printf("%T , The custmomer id ", customerId)
	VariantId, err := strconv.Atoi(variantId)

	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid Variant_id"})
		return
	}

	Quantity, err := strconv.Atoi(quantity)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid Quantity"})
		return
	}

	CustomerId, err := strconv.Atoi(customerId)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid CustomerId"})
		return
	}

	fmt.Println(VariantId, "  ", Quantity, " ", CustomerEmail, " ", CustomerId, " ", Address1, " ", City, " ", Province, " ", Country, " ", FinancialStatus)

	url := fmt.Sprintf("https://%s.myshopify.com/admin/api/%s/orders.json", config.ShopifyStore.ShopName, config.ShopifyStore.APIVersion)

	from := "venkatnaidu320@gmail.com"
	password := "dtqj unxa xerl uijk"

	smtpServer := "smtp.gmail.com"
	smtpPort := "587"
	to := []string{CustomerEmail}

	subject := "Subject: Order Creation\n"
	body := "Order Created Successfully"
	auth := smtp.PlainAuth("", from, password, smtpServer)
	message := []byte(subject + "\n" + body)

	orderJSON, err := json.Marshal(map[string]interface{}{
		"order": map[string]interface{}{
			"customer": map[string]interface{}{
				"id":    CustomerId,
				"email": CustomerEmail,
			},
			"line_items": []map[string]int{
				{
					"variant_id": VariantId,
					"quantity":   Quantity,
				},
			},
			"shipping_address": map[string]interface{}{
				"address1": Address1,
				"city":     City,
				"province": Province,
				"country":  Country,
			},
			"financial_status": FinancialStatus,
		},
	})

	if err != nil {
		log.Fatalf("Error marshaling order data: %v", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(orderJSON))
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	req.SetBasicAuth(config.ShopifyStore.APIKey, config.ShopifyStore.Password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
		message := []byte(subject + "\n" + "Order Not Created")

		smtp.SendMail(smtpServer+":"+smtpPort, auth, from, to, message)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		log.Printf("Failed to create order. Status: %v", resp.Status)
	} else {
		log.Println("Order created successfully!")
	}

	smtp.SendMail(smtpServer+":"+smtpPort, auth, from, to, message)

	c.JSON(http.StatusOK, gin.H{"message": "Order submitted successfully!"})

}
