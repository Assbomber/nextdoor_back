package auth

type RegisterRequest struct {
	// Required. User name. Max length=100, min length 1
	Username string `json:"username" binding:"required,max=100,min=1"`
	// Required. User email
	Email string `json:"email" binding:"required,email"`
	// Required. User password. Min 8 characters, max 100 characters
	Password string `json:"password" binding:"required,min=8,max=100"`
	// Required. OTP. min 111111, max 999999
	OTP int `json:"otp" binding:"required,min=111111,max=999999"`
}

type EmailVerificationRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type RegisterResponse struct {
	Token string `json:"token"`
}

type LoginRequest struct {
	// Required. User email
	EmailOrUsername string `json:"email" binding:"required,email"`
	// Required. User password. Min 8 characters, max 100 characters
	Password string `json:"password" binding:"required,min=8,max=100"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type ResetPasswordRequest struct {
	// Required. User email
	Email string `json:"email" binding:"required,email"`

	// Required. User password. Min 8 characters, max 100 characters
	Password string `json:"password" binding:"required,min=8,max=100"`

	// Required. OTP. min 111111, max 999999
	OTP int `json:"otp" binding:"required,min=111111,max=999999"`
}
