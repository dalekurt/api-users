package main

import (
	"fmt"
	"os"
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
	// Checking that the envronment variable is present or not for APP_PORT
	port, exists := os.LookupEnv("API_PORT")
	if !exists {
		fmt.Println("The environment variable API_PORT is not set")
	} else {
		fmt.Printf("Starting User API listening on port " + port)
	}
	// Init gin router
	router := gin.Default()

	// Its great to version your API's
	api := router.Group("/api")
	v1 := api.Group("/v1")
	users := v1.Group("/users")  
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
	users.POST("/signup", user.Signup)

	// Create the login endpoint
	users.POST("/login", user.Login)

	// Send reset link
	users.PUT("/reset-link", user.ResetLink)

	// Password reset
	users.PUT("/password-reset", user.PasswordReset)

	// Send verify link
	users.PUT("/verify-link", user.VerifyLink)

	// Verify account
	users.PUT("/verify-account", user.VerifyAccount)

	// Refresh token
	users.GET("/refresh", user.RefreshToken)

	// Handle error response when a route is not defined
	router.NoRoute(func(c *gin.Context) {
			// In gin this is how you return a JSON response
			c.JSON(404, gin.H{"message": "Not found"})
	})

	// TODO: Use env var for port
	// Init our server
	router.Run(":" + port)
}

