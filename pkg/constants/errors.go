package constants

import "errors"

var ErrUnexpectedSigningMethod = errors.New("unexpected JWT signing method found")
var ErrInvalidJWT = errors.New("the JWT is invalid")
var ErrNoSuchUser = errors.New("no user found")
var ErrWrongPassword = errors.New("incorrect password")
var ErrValidation = errors.New("invalid request body")
var ErrInvalidOTP = errors.New("otp not valid")
var ErrEmailAlreadyExist = errors.New("email already exists")
var ErrUsernameAlreadyExist = errors.New("username already exists")
var ErrNoAuthPresent = errors.New("you are not authenticated")
var ErrInvalidUserID = errors.New("invalid user id")
var ErrForbidden = errors.New("you are not authorized for this request")
