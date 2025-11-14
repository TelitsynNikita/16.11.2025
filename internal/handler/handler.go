package handler

import (
	"sync/atomic"
	"workmate/internal/service"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

type Handler struct {
	Service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h *Handler) InitRoutes(isShutDown *atomic.Bool) *fiber.App {
	router := fiber.New()

	router.Use(cors.New())

	// Мидлвара, которая нужна для предупреждения пользователей, что сервер находится в аварийном состоянии
	router.Use(func(c fiber.Ctx) error {
		if isShutDown.Load() {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"message": "service is shutting down",
			})
		}

		return c.Next()
	})

	message := router.Group("/link")
	message.Post("/check_by_urls", h.CheckLinksStatusByUrl)

	return router
}
