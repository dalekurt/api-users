package forms

// SignupUserCommand defines user form struct
type SignupUserCommand struct {
	// binding:"required" ensures that the field is provided
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// TODO: Allow user to login with username
// LoginUserCommand defines user login form struct
// type LoginEmailCommand struct {
// 	Email    string `json:"email" binding:"required"`
// 	Password string `json:"password" binding:"required"`
// }

type LoginUserCommand struct {
	Email string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ResendCommand defines resend email payload
type ResendCommand struct {
	// We only need the email to initialize an email sendout
	Email string `json:"email" binding:"required"`
}

// PasswordResetCommand defines user password reset form struct
type PasswordResetCommand struct {
	Password string `json:"password" binding:"required"`
	Confirm  string `json:"confirm" binding:"required"`
}
