package handler

import (
	"workmate/internal/model"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/jung-kurt/gofpdf"
)

func (h *Handler) CheckLinksStatusByUrl(c *fiber.Ctx) error {
	var body model.CheckLinksStatusByUrlRequest

	err := jsoniter.Unmarshal(c.Body(), &body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorMessage{
			Message:   err.Error(),
			ErrorCode: "marshal",
		})
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err = validate.Struct(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorMessage{
			Message:   err.Error(),
			ErrorCode: "validate",
		})
	}

	data, err := h.Service.PersistentURLService.CheckLinksStatusByUrl(body.Links)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorMessage{
			Message:   err.Error(),
			ErrorCode: "check_links_status_by_url",
		})
	}

	return c.Status(fiber.StatusOK).JSON(data)
}

func (h *Handler) CheckLinksStatusByID(c *fiber.Ctx) error {
	var body model.CheckLinksStatusByIDRequest

	err := jsoniter.Unmarshal(c.Body(), &body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorMessage{
			Message:   err.Error(),
			ErrorCode: "marshal",
		})
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err = validate.Struct(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorMessage{
			Message:   err.Error(),
			ErrorCode: "validate",
		})
	}

	links, err := h.Service.PersistentURLService.GetUrlByID(body.LinksList)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorMessage{
			Message:   err.Error(),
			ErrorCode: "check_links_status_by_id's",
		})
	}

	data, err := jsoniter.Marshal(links)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorMessage{
			Message:   err.Error(),
			ErrorCode: "marshal",
		})
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(190, 5, string(data), "0", "0", false)

	err = pdf.OutputFileAndClose("pdf_storage.pdf")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorMessage{
			Message:   err.Error(),
			ErrorCode: "output_file",
		})
	}

	return c.Download("./pdf_storage.pdf")
}
