package utils

import (
	"fmt"
	"time"

	"github.com/assbomber/myzone/pkg/constants"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

// Struct for JWT custom claims
type MyCustomClaims struct {
	UserID int64
	jwt.RegisteredClaims
}

// Generates JWT using user id, or else returns err if any.
func GenerateJWT(userID int64, jwtSecret string) (string, error) {
	claims := MyCustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(0, 6, 0)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

// Helps validate JWT using provided secret in JWT_SECRET environment variable.
// If Success, returns MyCustomClaims, else error
func ValidateJWT(tokenStr string, jwtSecret string) (*MyCustomClaims, error) {
	fmt.Println(jwtSecret)
	token, err := jwt.ParseWithClaims(tokenStr, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {

		// verifing if signing method is same
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, constants.ErrUnexpectedSigningMethod
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, errors.Wrap(constants.ErrInvalidJWT, "something")
	}

	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, constants.ErrInvalidJWT
	}
}
