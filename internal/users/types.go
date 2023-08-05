package users

import "time"

type OnboardingRequest struct {
	// User name
	Name      string `json:"name" binding:"required,min=1,max=100"`
	BirthDate string `json:"birthDate" binding:"required,date-check"`
	Gender    string `json:"gender" binding:"required,oneof=male female other"`
	Location  struct {
		Latitude  float64 `json:"latitude" binding:"required"`
		Longitude float64 `json:"longitude" binding:"required"`
	} `json:"location" binding:"required"`
}

type UserDetailsResponse struct {
	UserID    int64     `json:"userID"`
	Username  string    `json:"username"`
	Name      string    `json:"name"`
	Avatar    string    `json:"avatar"`
	BirthDate time.Time `json:"birthDate"`
	LastLogin time.Time `json:"lastLogin"`
	Gender    string    `json:"gender"`
	Latitude  float64   `json:"latitude" `
	Longitude float64   `json:"longitude"`
}
