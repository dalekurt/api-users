package main

import (
	// Log items to the terminal
	"log"

	// Import gin for route definition
	"github.com/gin-gonic/gin"
	// Import godotenv for .env variables
	"github.com/joho/godotenv"
	// Import our app controllers
	"github.com/dalekurt/api-users/controllers"
)

// init gets called before the main function
func init() {
	// Log error if .env file does not exist
	if err := godotenv.Load(); err != nil {
			log.Printf("No .env file found")
	}
}

func main() {
	// Init gin router
	router := gin.Default()

	// Its great to version your API's
	v1 := router.Group("/api/v1")
	{
			// Define the ping controller
			ping := new(controllers.PingController)
			// Define a GET request to call the Default
			// method in controllers/ping.go
			v1.GET("/ping", ping.Default)
	}

	// Define the user controller
	user := new(controllers.UserController)

	// Create the signup endpoint
	v1.POST("/signup", user.Signup)

	// Create the login endpoint
	v1.POST("/login", user.Login)

	// Handle error response when a route is not defined
	router.NoRoute(func(c *gin.Context) {
			// In gin this is how you return a JSON response
			c.JSON(404, gin.H{"message": "Not found"})
	})

	// Init our server
	router.Run(":3000")
}

