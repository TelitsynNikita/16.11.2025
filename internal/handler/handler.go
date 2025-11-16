package handler

import (
	"sync/atomic"
	"workmate/internal/service"

	"github.com/gofiber/fiber/v2"
)

var IsShutDown atomic.Bool

type Handler struct {
	Service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		Service: service,
	}
}

func (h *Handler) InitRoutes() *fiber.App {
	router := fiber.New()

	router.Use(func(c *fiber.Ctx) error {
		if IsShutDown.Load() {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"message": "service is shutting down",
			})
		}

		return c.Next()
	})

	message := router.Group("/link")
	message.Post("/check_by_urls", h.CheckLinksStatusByUrl)
	message.Post("/check_by_id", h.CheckLinksStatusByID)

	return router
}
