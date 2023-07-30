package auth

import (
	"net/http"

	"github.com/assbomber/myzone/pkg/logger"
	"github.com/assbomber/myzone/pkg/utils"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	log     *logger.Logger
	service Service
}

// Returns New Auth Handler
func NewHandler(log *logger.Logger, service Service) *Handler {
	return &Handler{
		log:     log,
		service: service,
	}
}

// Registers Auth Handler routes
func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	path := rg.Group("/auth")
	path.POST("/register", h.Register)
	path.POST("/login", h.Login)
}

// API handler for Register
func (h *Handler) Register(c *gin.Context) {
	var request RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.HandleErrorResponses(h.log, c, err)
		return
	}

	res, err := h.service.Register(c.Request.Context(), request)
	if err != nil {
		utils.HandleErrorResponses(h.log, c, err)
		return
	}

	c.JSON(http.StatusCreated, res)
}

// API handler for Login
func (h *Handler) Login(c *gin.Context) {
	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.HandleErrorResponses(h.log, c, err)
		return
	}

	res, err := h.service.Login(c.Request.Context(), request)
	if err != nil {
		utils.HandleErrorResponses(h.log, c, err)
		return
	}

	c.JSON(http.StatusCreated, res)
}
