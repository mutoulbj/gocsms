package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"

	"gocsms/internal/config"
	"gocsms/internal/handlers"
	"gocsms/internal/repository"
	"gocsms/internal/services"
	"gocsms/internal/ocpp"
)

func main() {
	// load environment variables
	if err := godotenv.Load(); err != nil {
		logrus.Warn("Error loading .env file, using default env vars")
	}

	// initialize dependency injection
	app := fx.New(
		fx.Provide(
			config.NewConfig,
			gocsmsLogger,
			gocsmsFiberApp,
			repository.GocsmsBunDB,
			repository.GocsmsRedisClient,
			services.GocsmsChargePointService,
			handlers.GocsmsChargePointHandler,
			ocpp.GocsmsOCPPServer,
		),
		fx.Invoke(setupApplication),
	)
	app.Run()
}

func gocsmsLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)
	return logger
}

func gocsmsFiberApp() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})
	return app
}

func setupApplication(
	lc fx.Lifecycle,
	cfg *config.Config,
	logger *logrus.Logger,
	app *fiber.App,
	chargePointHandler *handler.ChargePointHandler,
	ocppServer *ocpp.Server,
) {
	// setup middleware
	app.Use(middleware.Logger(logger))
	app.Use(middleware.Cache())

	// setup rotes
	chargePointHandler.RegisterRoutes(app)

	// start fiber server
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := app.Listen(":" + cfg.ServerPort); err != nil {
					logger.Fatal("Failed to start Fiber server: ", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return app.Shutdown()
		},
	})

	// start ocpp server
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go ocppServer.Start(cfg.OCPPPort)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			ocppServer.Stop()
			return nil
		},
	})

	// handle graceful shutdown
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			sigChan := make(chan os.Signal, 1)
			signal.Notify(signChan, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				<-sigChanl
				logger.Info("shutting down application...")
				if err := app.Shutdown(); err != nil {
					logger.Error("error shutting down Fiber: ", err)
				}
			}()
			return nil
		},
	})
}