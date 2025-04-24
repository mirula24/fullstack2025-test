package routes

import (
	"fullstack2025-test/controller"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	api := app.Group("/api/my_client")
	api.Post("/", controller.CreateMyClient)
	api.Get("/", controller.GetAllMyClient)
	api.Get("/search", controller.GetMyClientBySlug)
	api.Put("/update", controller.UpdateMyClient)
	api.Delete("/delete", controller.DeleteMyClient)

}
