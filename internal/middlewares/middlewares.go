package middlewares

import (
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
		return
	}

	claims, err := utils.ValidateJWT(bearerToken, m.jwtSecret)
	if err != nil {
		utils.HandleErrorResponses(m.log, c, err)
		return
	}

	c.Set(constants.USER_ID, claims.UserID)
	c.Next()
}
