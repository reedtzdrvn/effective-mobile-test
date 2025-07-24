package logger

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

var L *zap.Logger

func Init() {
	var err error
	L, err = zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
	}
}

func Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := c.Context().Time()
		err := c.Next()
		L.Info("request",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("latency", c.Context().Time().Sub(start)),
		)
		return err
	}
}
