package controller

import (
	"context"
	"database/constants"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var jwtKey = []byte("secret_key")

type LoginData struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func Login(client *mongo.Client, c *gin.Context) {
	credentials := LoginData{
		UserName: c.PostForm("user_name"),
		Password: c.PostForm("password"),
	}

	if credentials.UserName == "" || credentials.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_name and password are required"})
		return
	}

	collection := client.Database(constants.DatabaseName).Collection("registration_details")

	var user bson.M
	err := collection.FindOne(context.Background(), bson.M{"username": credentials.UserName}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error querying the database"})
		return
	}

	if user["password"] != credentials.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	expirationTime := time.Now().Add(time.Minute * 100)
	claims := &Claims{
		Username: credentials.UserName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.SetCookie("token", tokenString, int(expirationTime.Sub(time.Now()).Seconds()), "/", "", false, true)

	tokenStr, err := c.Cookie("token")
	if err != nil {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !tkn.Valid {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	if user["type"] == "client" {
		c.Redirect(http.StatusFound, "/clientpage")
	} else if user["type"] == "admin" {
		c.Redirect(http.StatusFound, "/adminpage")
	} else {
		c.Redirect(http.StatusFound, "/userpage")
	}
}
