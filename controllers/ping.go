package controllers

import (
    // Import the Gin library
    "github.com/gin-gonic/gin"
)

// PingController will hold the methods to the
type PingController struct{}

// Default controller handles returning the ping pong JSON response
func (h *PingController) Default(c *gin.Context) {
    c.JSON(200, gin.H{"message": "pong"})
}


