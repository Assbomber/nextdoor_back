package users

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
