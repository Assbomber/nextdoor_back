package auth

type RegisterRequest struct {
	Name     string `json:"name" binding:"required,max=100,min=1"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=100"`
}

type RegisterResponse struct {
	Token string `json:"token"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=100"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
