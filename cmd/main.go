package main

import (
	"github.com/gofiber/fiber/v3"
	"github.com/sirupsen/logrus"
)

func main() {
	app := fiber.New()

	// Добавить хэндлеры
	app.Server().Handler = nil

	if err := app.Listen(":8080"); err != nil {
		logrus.Fatalf("error starting server: %v", err)
	}
}
