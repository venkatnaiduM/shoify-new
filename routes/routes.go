package routes

import (
	"database/controller"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func Routing(client *mongo.Client) {

	r := gin.Default()
	r.Use(cors.Default())
	r.LoadHTMLFiles("form.html", "register.html", "login.html", "home.html", "admin.html", "client.html", "user.html", "cart.html", "order.html")
	r.GET("/homepage", controller.Home)
	r.GET("/orderform", controller.OrderHandle)
	r.GET("/adminpage", controller.Admin)
	r.GET("/clientpage", controller.Client)
	r.GET("/userpage", controller.User)
	r.GET("/", controller.ServeForm)
	r.GET("/register", controller.RegistraionForm)
	r.GET("/loginform", controller.LoginForm)
	r.GET("/cartform", controller.Cart)
	r.POST("/login", func(c *gin.Context) {
		controller.Login(client, c)
	})
	r.POST("/submit", func(c *gin.Context) {
		controller.SubmitHandler(client, c)
	})
	r.POST("/registrationdetails", func(c *gin.Context) {
		controller.RegistrationHandler(client, c)
	})
	r.POST("/delete", func(c *gin.Context) {
		controller.DeleteData(client, c)
	})
	r.POST("/update", func(c *gin.Context) {
		controller.UpdateDetails(client, c)
	})
	r.POST("/submitorder", controller.HandleOrderSubmission)
	r.POST("/addtocart", controller.AddToCart)
	r.POST("/removefromcart", controller.RemoveFromCart)
	r.GET("/cart", controller.GetCart)
	r.GET("/products", controller.ProductDetails)

	r.GET("/checkouts", controller.CheckOutDetails)
	r.GET("/pricerules", controller.PriceRules)

	r.GET("/discountcodes", controller.DiscountCodes)

	r.GET("/getcartdetails", controller.GetCartDetails)
	r.GET("/deletecartdetails", controller.DeleteCartDetails)
	r.GET("/cartlinedetails", controller.CartLinesAdd)

	r.GET("/orders", controller.OrderDetails)

	r.GET("/customers", controller.CustomerDetails)

	if err := r.Run(":9090"); err != nil {
		log.Fatalf("Failed to run the server: %v", err)
	}

}
