package controllers

import (
    "github.com/gin-gonic/gin"
    "github.com/dalekurt/api-users/forms"
    "github.com/dalekurt/api-users/models"
)

// Import the userModel from the models
var userModel = new(models.UserModel)

// UserController defines the user controller methods
type UserController struct{}

// Signup controller handles registering a user
func (u *UserController) Signup(c *gin.Context) {
    var data forms.SignupUserCommand

    // Bind the data from the request body to the SignupUserCommand Struct
    // Also check if all fields are provided
    if c.BindJSON(&data) != nil {
        // specified response
        c.JSON(406, gin.H{"message": "A required field is missing"})
        // abort the request
        c.Abort()
        // return nothing
        return
    }

    /*
        You can add your validation logic
        here such as email

        if regexMethodChecker(data.Email) {
            c.JSON(400, gin.H{"message": "Email is invalid"})
            c.Abort()
            return
        }
    */

    err := userModel.Signup(data)

    // Check if there was an error when saving user
    if err != nil {
        c.JSON(400, gin.H{"message": "There was a problem creating your account"})
        c.Abort()
        return
    }

    c.JSON(201, gin.H{"message": "Your account has been created"})
}