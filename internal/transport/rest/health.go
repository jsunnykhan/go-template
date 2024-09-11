package transport

import (
	"name/internal/models"

	"github.com/gofiber/fiber/v2"
)

func GetHealthStatus(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).JSON(models.APIResponse{
		Code:    fiber.StatusOK,
		Message: "Server is running",
	})
}
