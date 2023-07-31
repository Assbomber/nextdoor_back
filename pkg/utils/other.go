package utils

// redis Key for Otp verification
func GetOTPRedisKey(email string) string {
	return "MYZONE-VERIFY-" + email
}
