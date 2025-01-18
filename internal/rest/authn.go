package rest

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/muhrifqii/tuskar/domain"
	"github.com/muhrifqii/tuskar/internal/config"
	"github.com/muhrifqii/tuskar/internal/rest/rest_utils"
	"go.uber.org/zap"
)

type (
	AuthnService interface {
		Login(ctx context.Context, req domain.AuthnRequest) (domain.AuthnResponse, error)
		Logout(ctx context.Context) error
	}

	AuthnHandler struct {
		authnService AuthnService
		validator    *validator.Validate
		log          *zap.Logger
		conf         *config.JwtConfig
	}
)

func NewAuthnHandler(router fiber.Router, svc AuthnService, params rest_utils.HandlerParams, jwtConf config.JwtConfig) {
	handler := &AuthnHandler{
		authnService: svc,
		validator:    params.Validator,
		log:          params.Logger,
		conf:         &jwtConf,
	}

	router.Post("/authenticate", handler.Login)

}

func NewProtectedAuthnHandler(router fiber.Router, svc AuthnService, params rest_utils.HandlerParams) {
	handler := &AuthnHandler{
		authnService: svc,
		validator:    params.Validator,
		log:          params.Logger,
	}

	protectedAuthnRoute := router.Group("/authenticate")
	protectedAuthnRoute.Put("", handler.Refresh)
}

func (h *AuthnHandler) Login(c *fiber.Ctx) error {
	var req domain.AuthnRequest
	if err := c.BodyParser(&req); err != nil {
		return rest_utils.NewApiErrorResponse(fiber.StatusBadRequest, err.Error())
	}
	if err := h.validator.Struct(&req); err != nil {
		return err
	}
	response, err := h.authnService.Login(c.Context(), req)
	if err != nil {
		return err
	}
	c.Cookie(&fiber.Cookie{
		Name:     h.conf.CookieName,
		Value:    response.RefreshToken,
		Expires:  response.RefreshTokenExpiresAt,
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Lax",
		Path:     c.Path(),
	})
	return rest_utils.ReturnOkResponse(c, response)
}

func (h *AuthnHandler) Refresh(c *fiber.Ctx) error {
	return nil
}
