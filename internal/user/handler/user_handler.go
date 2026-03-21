package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/manav1011/ikatva-be/internal/user/model"
	"github.com/manav1011/ikatva-be/internal/user/service"
	"github.com/manav1011/ikatva-be/pkg/utils"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

// Login authenticates a user and returns access and refresh tokens.
// @Summary      User login
// @Description  Validates email and password; returns JWT access and refresh tokens.
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        body  body      model.LoginRequest  true  "Email and password"
// @Success      200   {object}  model.LoginSuccessEnvelope
// @Failure      400   {object}  utils.APIResponse
// @Failure      401   {object}  utils.APIResponse
// @Failure      500   {object}  utils.APIResponse
// @Router       /users/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, "invalid request body", http.StatusBadRequest)
		return
	}

	data, err := h.svc.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			utils.Error(c, "invalid credentials", http.StatusUnauthorized)
			return
		}
		utils.Error(c, "internal server error", http.StatusInternalServerError)
		return
	}

	utils.Success(c, data)
}
