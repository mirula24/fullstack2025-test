package main

import (
	"fullstack2025-test/database"
	"fullstack2025-test/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}
	database.ConnectDB()
	database.ConnectRedis()

	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Terjadi kesalahan pada server",
				"error":   err.Error(),
			})
		},
	})
	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
