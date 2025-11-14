package handler

import (
	"github.com/gofiber/fiber/v3"
)

func (h *Handler) CheckLinksStatusByUrl(c fiber.Ctx) error {
	id, err := h.Service.URLService.GetUrlByID(0)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"url": id,
	})
}

func (h *Handler) CheckLinksStatusByID(c fiber.Ctx) error {
	id, err := h.Service.URLService.GetUrlByID(0)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"url": id,
	})
}
