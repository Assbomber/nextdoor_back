package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/assbomber/myzone/pkg/constants"
	"github.com/assbomber/myzone/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

type Response struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type StackTracer interface {
	StackTrace() errors.StackTrace
}

// Initializes custom gin binding validations
func InitGinCustomValidations() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {

		v.RegisterValidation("date-check", func(fl validator.FieldLevel) bool {
			val := fl.Field().String()
			if !ok {
				return false
			}

			// validating date
			_, err := time.Parse("2006-01-02", val)
			if err != nil {
				return false
			}
			return true
		})
	}
}

// Helper that handles writing error responses to gin context. Does not abort.
func HandleErrorResponses(log *logger.Logger, c *gin.Context, errr error, extra ...string) {
	err := errors.Cause(errr)
	// *Gin binding validations. Checks whether is gin binding errors -------------------------
	// if validation err, extracting actual err
	sliceErrs := make([]BindingErrorMsg, 0)
	if _, ok := err.(validator.ValidationErrors); ok {
		errs := getBindingErrors(err)
		sliceErrs = append(sliceErrs, errs...)
		err = constants.ErrValidation
	} else if er, ok := err.(*json.UnmarshalTypeError); ok {
		fmt.Println(er.Field, er.Type, er.Value)
		err = constants.ErrValidation
		sliceErrs = append(sliceErrs, BindingErrorMsg{Field: er.Field, Message: "Must be of type " + er.Type.String()})
	}

	// if slice validation err, extracting actual err
	if errrs, ok := err.(binding.SliceValidationError); ok {
		for _, err := range errrs {
			errs := getBindingErrors(err)
			sliceErrs = append(sliceErrs, errs...)
		}
		err = constants.ErrValidation
	}
	// *----------------------------------------------------------------------------------------

	// detecting err type and preparing reponses
	switch err {

	case constants.ErrValidation:
		if len(sliceErrs) > 0 {
			c.JSON(http.StatusBadRequest, Response{Message: err.Error(), Data: sliceErrs})
		} else {
			c.JSON(http.StatusBadRequest, Response{Message: err.Error() + ", " + strings.Join(extra, ",")})
		}
	case jwt.ErrTokenMalformed:
		c.JSON(http.StatusUnauthorized, Response{Message: err.Error()})
	case constants.ErrEmailAlreadyExist:
		c.JSON(http.StatusConflict, Response{Message: err.Error()})
	case constants.ErrInvalidOTP:
		c.JSON(http.StatusUnauthorized, Response{Message: err.Error()})
	case constants.ErrUnexpectedSigningMethod:
		c.JSON(http.StatusUnauthorized, Response{Message: err.Error()})
	case constants.ErrInvalidJWT:
		c.JSON(http.StatusUnauthorized, Response{Message: err.Error()})
	case constants.ErrInvalidUserID:
		c.JSON(http.StatusBadRequest, Response{Message: err.Error()})
	case constants.ErrNoSuchUser:
		c.JSON(http.StatusNotFound, Response{Message: err.Error()})
	case constants.ErrWrongPassword:
		c.JSON(http.StatusUnauthorized, Response{Message: err.Error()})
	case constants.ErrUsernameAlreadyExist:
		c.JSON(http.StatusConflict, Response{Message: err.Error()})
	case constants.ErrNoAuthPresent:
		c.JSON(http.StatusUnauthorized, Response{Message: err.Error()})
	case constants.ErrForbidden:
		c.JSON(http.StatusForbidden, Response{Message: err.Error()})
	default:
		// printing stack trace
		var errs string = errr.Error() + "\n"
		if err, ok := errr.(StackTracer); ok {
			for _, f := range err.StackTrace() {
				errs += fmt.Sprintf("%+s:%d\n", f, f)
			}
		}
		log.Error(errs, nil)
		c.JSON(http.StatusInternalServerError, Response{Message: "Oops! its not you, its us. Please try again later."})
	}
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "field is required"
	case "oneof":
		return "must be one of: " + fe.Param()
	case "min":
		return "must be greater than or equal to " + fe.Param()
	case "max":
		return "must be lesserr than or equal to " + fe.Param()
	case "email":
		return "invalid email " + fe.Param()
	case "date-check":
		return "invalid date, format required is YYYY-MM-DD" + fe.Param()
	}

	return "Unknown error"
}

type BindingErrorMsg struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func getBindingErrors(err error) []BindingErrorMsg {

	errs := make([]BindingErrorMsg, 0)
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, fe := range ve {
			errs = append(errs, BindingErrorMsg{Field: fe.Field(), Message: getErrorMsg(fe)})
		}
	}
	return errs
}
