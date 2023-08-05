package utils

import "time"

// redis Key for Otp verification
func GetOTPRedisKey(email string) string {
	return "MYZONE-VERIFY-" + email
}

// Returns IST location for time
func GetISTLocation() *time.Location {
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		panic(" Error loading Kolkata location")
	}
	return loc
}
