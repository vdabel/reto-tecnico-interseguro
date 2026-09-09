package main

import (
	"fmt"
	"log"

	"api-go/config"
	"api-go/handlers"
	"api-go/middleware"
	"api-go/services"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	cfg := config.LoadConfig()

	app := fiber.New(fiber.Config{
		AppName:      "Interseguro - Go Matrix Processing API",
		ServerHeader: "Fiber",
	})

	// Middlewares globales
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, OPTIONS",
	}))

	// Inicializar capas de servicio y controladores
	matrixService := services.NewMatrixService()
	authHandler := handlers.NewAuthHandler(cfg)
	matrixHandler := handlers.NewMatrixHandler(cfg, matrixService)

	// Rutas Públicas
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"service": "api-go",
			"version": "1.0.0",
		})
	})

	app.Post("/api/auth/token", authHandler.GenerateToken)

	// Rutas Protegidas por JWT
	api := app.Group("/api", middleware.ProtectedMiddleware(cfg.JWTSecret))
	api.Post("/matrix/process", matrixHandler.ProcessMatrix)

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("🚀 API en Go escuchando en el puerto %s...", cfg.Port)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("Error al iniciar el servidor Fiber: %v", err)
	}
}
