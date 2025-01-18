package middleware

import (
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/muhrifqii/tuskar/internal/config"
	"github.com/muhrifqii/tuskar/internal/rest/rest_utils"
)

func RequireAuthn(conf config.JwtConfig) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: []byte(conf.Secret),
		},
		ErrorHandler: rest_utils.JwtErrorResponseHandler,
	})
}
