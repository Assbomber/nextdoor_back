package utils

import (
	"fmt"
	"net/http"

	"github.com/assbomber/myzone/pkg/constants"
	"github.com/assbomber/myzone/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"
)

type ErrResponse struct {
	Message string `json:"message"`
}
type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type StackTracer interface {
	StackTrace() errors.StackTrace
}

// Helper that handles writing error responses to gin context. Does not abort.
func HandleErrorResponses(log *logger.Logger, c *gin.Context, err error, extra ...string) {
	// *Gin binding validations. Checks whether is gin binding errors -------------------------
	// if validation err, extracting actual err
	sliceErrs := make([]BindingErrorMsg, 0)
	if _, ok := err.(validator.ValidationErrors); ok {
		errs := getBindingErrors(err)
		sliceErrs = append(sliceErrs, errs...)
		err = constants.ErrValidation
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
			c.JSON(http.StatusBadRequest, sliceErrs)
		} else {
			c.JSON(http.StatusBadRequest, ErrResponse{Message: err.Error()})
		}
	case constants.ErrEmailAlreadyExist:
		c.JSON(http.StatusBadRequest, ErrResponse{Message: err.Error()})
	case constants.ErrInvalidOTP:
		c.JSON(http.StatusUnauthorized, ErrResponse{Message: err.Error()})
	case constants.ErrUnexpectedSigningMethod:
		c.JSON(http.StatusUnauthorized, ErrResponse{Message: err.Error()})
	case constants.ErrInvalidJWT:
		c.JSON(http.StatusUnauthorized, ErrResponse{Message: err.Error()})
	case constants.ErrNoSuchUser:
		c.JSON(http.StatusNotFound, ErrResponse{Message: err.Error()})
	case constants.ErrWrongPassword:
		c.JSON(http.StatusUnauthorized, ErrResponse{Message: err.Error()})
	case constants.ErrUsernameAlreadyExist:
		c.JSON(http.StatusBadRequest, ErrResponse{Message: err.Error()})
	default:
		// printing stack trace
		var errs string = err.Error() + "\n"
		if err, ok := err.(StackTracer); ok {
			for _, f := range err.StackTrace() {
				errs += fmt.Sprintf("%+s:%d\n", f, f)
			}
		}
		log.Error(errs, nil)
		c.JSON(http.StatusInternalServerError, ErrResponse{Message: "Oops! its not you, its us. Please try again later."})
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
