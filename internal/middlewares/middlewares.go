package middlewares

import (
	"strconv"

	"github.com/assbomber/myzone/pkg/constants"
	"github.com/assbomber/myzone/pkg/logger"
	"github.com/assbomber/myzone/pkg/utils"
	"github.com/gin-gonic/gin"
)

type Middleware struct {
	log       *logger.Logger
	jwtSecret string
}

func New(log *logger.Logger, jwtSecret string) *Middleware {
	return &Middleware{
		log:       log,
		jwtSecret: jwtSecret,
	}
}

// Only allows user to proceed if Authorization header is present and token is valid
func (m *Middleware) MustBeLoggedIn(c *gin.Context) {
	bearerToken := c.GetHeader("Authorization")
	if bearerToken == "" {
		utils.HandleErrorResponses(m.log, c, constants.ErrNoAuthPresent)
		c.Abort()
		return
	}

	claims, err := utils.ValidateJWT(bearerToken, m.jwtSecret)
	if err != nil {
		utils.HandleErrorResponses(m.log, c, err)
		c.Abort()
		return
	}

	c.Set(constants.USER_ID, claims.UserID)
	c.Next()
}

// Only allows user to proceed if self request or admin request
func (m *Middleware) MustBeSelfOrAdmin(c *gin.Context) {
	bearerToken := c.GetHeader("Authorization")
	if bearerToken == "" {
		utils.HandleErrorResponses(m.log, c, constants.ErrNoAuthPresent)
		c.Abort()
		return
	}

	claims, err := utils.ValidateJWT(bearerToken, m.jwtSecret)
	if err != nil {
		utils.HandleErrorResponses(m.log, c, err)
		c.Abort()
		return
	}

	paramUser, err := strconv.ParseInt(c.Param(constants.USER_ID), 10, 64)
	if err != nil {
		utils.HandleErrorResponses(m.log, c, constants.ErrInvalidUserID)
		c.Abort()
		return
	}
	// Not same user, also not an admin
	if paramUser != claims.UserID && !claims.IsAdmin {
		utils.HandleErrorResponses(m.log, c, constants.ErrForbidden)
	}

	c.Set(constants.USER_ID, claims.UserID)
	c.Next()
}
