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
	path.POST("/email_verification", h.SendVerificationEmail)
	path.POST("/reset", h.ResetPassword)
}

// Register godoc
// @Summary      API for new user registration.
// @Tags         auth
// @Produce      json
// @Param        body   body      RegisterRequest  true  "Request body"
// @success 	 200 {object} RegisterResponse
// @Failure      400  {object}  utils.ErrResponse
// @Failure      401  {object}  utils.ErrResponse
// @Failure      403  {object}  utils.ErrResponse
// @Failure      500  {object}  utils.ErrResponse
// @Router       /auth/register [post]
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

// Login godoc
// @Summary      API for user login.
// @Tags         auth
// @Produce      json
// @Param        body   body      LoginRequest  true  "Request body"
// @success 	 200 {object} LoginResponse
// @Failure      400  {object}  utils.ErrResponse
// @Failure      401  {object}  utils.ErrResponse
// @Failure      403  {object}  utils.ErrResponse
// @Failure      500  {object}  utils.ErrResponse
// @Router       /auth/login [post]
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

// SendVerificationEmail godoc
// @Summary      Sends verification email.
// @Tags         auth
// @Produce      json
// @Param        body   body      EmailVerificationRequest  true  "Request body"
// @success 	 200 {object} utils.ErrResponse
// @Failure      400  {object}  utils.ErrResponse
// @Failure      401  {object}  utils.ErrResponse
// @Failure      403  {object}  utils.ErrResponse
// @Failure      500  {object}  utils.ErrResponse
// @Router       /auth/email_verification [post]
func (h *Handler) SendVerificationEmail(c *gin.Context) {
	var request EmailVerificationRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.HandleErrorResponses(h.log, c, err)
		return
	}

	err := h.service.SendVerificationEmail(c.Request.Context(), request.Email)
	if err != nil {
		utils.HandleErrorResponses(h.log, c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "verification email sent"})
}

// ResetPassword godoc
// @Summary      Resets user password.
// @Tags         auth
// @Produce      json
// @Param        body   body      ResetPasswordRequest  true  "Request body"
// @success 	 200 {object} utils.ErrResponse
// @Failure      400  {object}  utils.ErrResponse
// @Failure      401  {object}  utils.ErrResponse
// @Failure      403  {object}  utils.ErrResponse
// @Failure      500  {object}  utils.ErrResponse
// @Router       /auth/reset [post]
func (h *Handler) ResetPassword(c *gin.Context) {
	var request ResetPasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		utils.HandleErrorResponses(h.log, c, err)
		return
	}

	err := h.service.ResetPassword(c.Request.Context(), request)
	if err != nil {
		utils.HandleErrorResponses(h.log, c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password reset successful"})
}
