package users

import (
	"net/http"

	"github.com/assbomber/myzone/internal/middlewares"
	"github.com/assbomber/myzone/pkg/constants"
	"github.com/assbomber/myzone/pkg/logger"
	"github.com/assbomber/myzone/pkg/utils"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	log         *logger.Logger
	middlewares *middlewares.Middleware
	service     Service
}

// Returns New Auth Handler
func NewHandler(log *logger.Logger, middlewares *middlewares.Middleware, service Service) *Handler {
	return &Handler{
		log:         log,
		middlewares: middlewares,
		service:     service,
	}
}

// Registers Auth Handler routes
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	path := rg.Group("/users")
	path.POST("/onboarding", h.middlewares.MustBeLoggedIn, h.CreateOnboardingDetails)
}

// CreateOnboardingDetails godoc
// @Summary      Api thats create basic onboarding details for a user
// @Tags         users
// @Produce      json
// @Param        body   body      OnboardingRequest  true  "Request body"
// @success 	 200 {object} utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      401  {object}  utils.Response
// @Failure      403  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /users/onboarding [post]
// @Security ApiKeyAuth
func (h *Handler) CreateOnboardingDetails(c *gin.Context) {
	var request OnboardingRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.HandleErrorResponses(h.log, c, err)
		return
	}

	userID := c.GetInt64(constants.USER_ID)

	err := h.service.CreateOnboardingDetails(c.Request.Context(), userID, request)
	if err != nil {
		utils.HandleErrorResponses(h.log, c, err)
		return
	}

	c.JSON(http.StatusCreated, utils.Response{Message: "Onboarding details saved successfully"})
}
