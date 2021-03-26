package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/dalekurt/api-users/forms"
	"github.com/dalekurt/api-users/models"
	"github.com/dalekurt/api-users/helpers"
	"github.com/dalekurt/api-users/services"
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

	result, _ := userModel.GetUserByEmail(data.Email)

	// If there happens to be a result respond with a 
	// descriptive mesage
	if result.Email != "" {
			c.JSON(403, gin.H{"message": "This email has already registered"})
			c.Abort()
			return
	}

	err := userModel.Signup(data)

	// Check if there was an error when saving user
	if err != nil {
			c.JSON(400, gin.H{"message": "There was a problem creating your account"})
			c.Abort()
			return
	}

	c.JSON(201, gin.H{"message": "Your account has been created"})
	
	// Generate token to hold users details
	resetToken, _ := services.GenerateNonAuthToken(data.Email)

	// link to be verify account
	link := "http://localhost:5000/api/v1/verify-account?verify_token=" + resetToken
	// Define email body
	body := "Here is your reset <a href='" + link + "'>link</a>"
	html := "<strong>" + body + "</strong>"

	// initialize email send out
	email := services.SendMail("Verify Account", body, data.Email, html, data.Username)

	// If email fails while sending
	if !email {
			c.JSON(500, gin.H{"message": "An issue occured sending you an email"})
			c.Abort()
			return
	}
}

// Login allows a user to login a user and get
// access token
func (u *UserController) Login(c *gin.Context) {
	var data forms.LoginUserCommand

	// Bind the request body data to var data and check if all details are provided
	if c.BindJSON(&data) != nil {
			c.JSON(406, gin.H{"message": "Invalid or missing credentials"})
			c.Abort()
			return
	}

	result, err := userModel.GetUserByEmail(data.Email)

	if result.Email == "" {
			c.JSON(404, gin.H{"message": "Authentication Error, User Not found"})
			c.Abort()
			return
	}

	if err != nil {
			c.JSON(400, gin.H{"message": "Login error"})
			c.Abort()
			return
	}

	// Get the hashed password from the saved document
	hashedPassword := []byte(result.Password)
	// Get the password provided in the request.body
	password := []byte(data.Password)

	err = helpers.PasswordCompare(password, hashedPassword)

	if err != nil {
			c.JSON(403, gin.H{"message": "Invalid user credentials"})
			c.Abort()
			return
	}

	jwtToken, refreshToken, err2 := services.GenerateToken(data.Email)

	// If we fail to generate token for access
	if err2 != nil {
			c.JSON(403, gin.H{"message": "There was a problem logging you in, try again later"})
			c.Abort()
			return
	}

	c.JSON(200, gin.H{"message": "Log in success", "token": jwtToken, "refresh_token": refreshToken})
}

// ResetLink handles resending email to user to reset link
func (u *UserController) ResetLink(c *gin.Context) {
	// Defined schema for the request body
	var data forms.ResendCommand

	// Ensure the user provides all values from the request.body
	if (c.BindJSON(&data)) != nil {
			// Return 400 status if they don't provide the email
			c.JSON(400, gin.H{"message": "Provided all fields"})
			c.Abort()
			return
	}

	// Fetch the account from the database based on the email
	// provided
	result, err := userModel.GetUserByEmail(data.Email)

	// Return 404 status if an account was not found
	if result.Email == "" {
			c.JSON(404, gin.H{"message": "User account was not found"})
			c.Abort()
			return
	}

	// Return 500 status if something went wrong while fetching
	// account
	if err != nil {
			c.JSON(500, gin.H{"message": "Error encountered while fetching account, please login again"})
			c.Abort()
			return
	}

	// Generate the token that will be used to reset the password
	resetToken, _ := services.GenerateNonAuthToken(result.Email)

	// The link to be clicked in order to perform a password reset
	link := "http://localhost:5000/api/v1/password-reset?reset_token=" + resetToken
	// Define the body of the email
	body := "Here is your reset <a href='" + link + "'>link</a>"
	html := "<strong>" + body + "</strong>"

	// Initialize email sendout
	email := services.SendMail("Reset Password", body, result.Email, html, result.Username)

	// If email was sent, return 200 status code
	if email == true {
		c.JSON(200, gin.H{"messsage": "Please check your email"})
		c.Abort()
		return
		// Return 500 status when something wrong happened
	} else {
		c.JSON(500, gin.H{"message": "An issue occured sending you an email"})
		c.Abort()
		return
	}
	}
	
// PasswordReset handles user password request
func (u *UserController) PasswordReset(c *gin.Context) {
	var data forms.PasswordResetCommand

	// Ensure they provide data based on the schema
	if c.BindJSON(&data) != nil {
		c.JSON(406, gin.H{"message": "Provide required fields"})
		c.Abort()
		return
	}

	// Ensures that the password provided matches the confirm
	if data.Password != data.Confirm {
		c.JSON(400, gin.H{"message": "Passwords do not match"})
		c.Abort()
		return
	}

	// Get token from link query sent to your email
	resetToken, _ := c.GetQuery("reset_token")

	// Decode the token
	userID, _ := services.DecodeNonAuthToken(resetToken)

	// Fetch the user
	result, err := userModel.GetUserByEmail(userID)

	if err != nil {
		// Return response when we get an error while fetching user
		c.JSON(500, gin.H{"message": "Something went wrong. Please try again."})
		c.Abort()
		return
	}
	// Check if account exists
	if result.Email == "" {
		c.JSON(404, gin.H{"message": "Your user account was not found."})
		c.Abort()
		return
	}

	// Hash the new password
	newHashedPassword := helpers.GeneratePasswordHash([]byte(data.Password))

	// Update user account
	_err := userModel.UpdateUserPass(userID, newHashedPassword)

	if _err != nil {
		// Return response if we are not able to update user password
		c.JSON(500, gin.H{"message": "Something went wrong while updating your password. Please try again."})
		c.Abort()
		return
	}

	c.JSON(201, gin.H{"message": "Your password has been updated. Please login."})
	c.Abort()
	return
}
	
// VerifyLink handles resending email to user to reset link
func (u *UserController) VerifyLink(c *gin.Context) {
	var data forms.ResendCommand

	// Ensure they provide all relevant fields in the request body
	if (c.BindJSON(&data)) != nil {
		c.JSON(400, gin.H{"message": "Provided all fields"})
		c.Abort()
		return
	}

	// Fetch account from database
	result, err := userModel.GetUserByEmail(data.Email)

	// Check if account exist return 404 if not
	if result.Email == "" {
		c.JSON(404, gin.H{"message": "User account was not found"})
		c.Abort()
		return
	}

	if err != nil {
		c.JSON(500, gin.H{"message": "Something wrong happened, try again later"})
		c.Abort()
		return
	}

	// Generate token to hold user details
	resetToken, _ := services.GenerateNonAuthToken(result.Email)

	// Define email body
	link := "http://localhost:5000/api/v1/verify-account?verify_token=" + resetToken
	body := "Here is your reset <a href='" + link + "'>link</a>"
	html := "<strong>" + body + "</strong>"

	// Initialize email sendout
	email := services.SendMail("Verify Account", body, result.Email, html, result.Username)

	// If email send 200 status code
	if email == true {
		c.JSON(200, gin.H{"messsage": "Check mail"})
		c.Abort()
		return
	} else {
		c.JSON(500, gin.H{"message": "An issue occured sending you an email"})
		c.Abort()
		return
	}
}
	
// VerifyAccount handles user password request
func (u *UserController) VerifyAccount(c *gin.Context) {
	// Get token from link query
	verifyToken, _ := c.GetQuery("verify_token")

	// Decode verify token
	userID, _ := services.DecodeNonAuthToken(verifyToken)

	// Fetch user based on details from decoded token
	result, err := userModel.GetUserByEmail(userID)

	if err != nil {
		// Return response when we get an error while fetching user
		c.JSON(500, gin.H{"message": "Something wrong happened, try again later"})
		c.Abort()
		return
	}

	if result.Email == "" {
		c.JSON(404, gin.H{"message": "User account was not found"})
		c.Abort()
		return
	}

	// Update user account
	_err := userModel.VerifyAccount(userID)

	if _err != nil {
		// Return response if we are not able to update user password
		c.JSON(500, gin.H{"message": "Something happened while verifying you account, try again"})
		c.Abort()
		return
	}

	c.JSON(201, gin.H{"message": "Account verified, log in"})
}

// RefreshToken handles refresh token
func (u *UserController) RefreshToken(c *gin.Context) {
	refreshToken := c.Request.Header["Refreshtoken"]

	if refreshToken == nil {
		c.JSON(403, gin.H{"message": "No refresh token provided"})
		c.Abort()
		return
	}

	email, err := services.DecodeRefreshToken(refreshToken[0])

	if err != nil {
		c.JSON(500, gin.H{"message": "Problem refreshing your session"})
		c.Abort()
		return
	}

	// Create new token
	accessToken, _refreshToken, _err := services.GenerateToken(email)

	if _err != nil {
		c.JSON(500, gin.H{"message": "There was a problem creating a new session"})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{"message": "Your login was successful", "token": accessToken, "refresh_token": _refreshToken})
}
	