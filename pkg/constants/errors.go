package constants

import "errors"

var ErrUnexpectedSigningMethod = errors.New("unexpected JWT signing method found")
var ErrInvalidJWT = errors.New("the JWT is invalid")
var ErrNoSuchUser = errors.New("no user found")
var ErrWrongPassword = errors.New("wrong password")
var ErrValidation = errors.New("invalid request body")
