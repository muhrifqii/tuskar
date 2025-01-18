package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/muhrifqii/tuskar/internal/config"
)

func RequestID(conf config.ApiConfig) fiber.Handler {
	return requestid.New(requestid.Config{
		Header:     conf.HeaderRequestID,
		Generator:  utils.UUIDv4,
		ContextKey: "requestId",
	})
}
