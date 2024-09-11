package api

import (
	transport "name/internal/transport/rest"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func InitializeNewServer() *fiber.App {
	app := fiber.New()
	api := app.Group("/api/v1")
	api.Use(logger.New(loggerConfig()))
	configureCors(api)
	registerHandlers(api)
	return app
}

func loggerConfig() logger.Config {

	return logger.Config{
		Format:        "Form ${ip} requested for ${method} ${path} at ${time} -- ${status} and latency was ${latency}\n",
		TimeFormat:    "02-Jan-2006:15:04:05",
		TimeZone:      "Local",
		DisableColors: false,
	}
}

func configureCors(app fiber.Router) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET, POST, PUT, DELETE",
	}))
}

func registerHandlers(router fiber.Router) {
	router.Get("health", transport.GetHealthStatus)
}
