package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/effective-mobile/subscriptions/internal/config"
	"github.com/effective-mobile/subscriptions/internal/db"
	"github.com/effective-mobile/subscriptions/internal/logger"
	"github.com/effective-mobile/subscriptions/internal/repository"
	httpapi "github.com/effective-mobile/subscriptions/internal/transport/http"
	"github.com/effective-mobile/subscriptions/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config load error: %v", err)
	}
	logger.Init()
	lg := logger.L
	pool, err := db.NewPostgresPool(cfg)
	if err != nil {
		lg.Fatal("db connect error", zap.Error(err))
	}
	defer pool.Close()
	repo := repository.NewSubscriptionPostgres(pool)
	uc := usecase.NewSubscriptionUseCase(repo)
	handler := httpapi.NewSubscriptionHandler(uc)
	app := fiber.New()
	app.Use(logger.Middleware())
	handler.RegisterRoutes(app)

	// Graceful shutdown
	go func() {
		if err := app.Listen(":" + cfg.AppPort); err != nil {
			lg.Fatal("fiber listen error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err := app.Shutdown(); err != nil {
		lg.Error("fiber shutdown error", zap.Error(err))
	}
	lg.Info("server stopped")
}
