package router

import (
	"github.com/farid141/go-rest-api/controller"
	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Get("/users", controller.GetUsers)
	api.Post("/users", controller.CreateUser)
}

func SetupPublicRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Post("/login", controller.Login)
	api.Post("/logout", controller.Logout)
}
