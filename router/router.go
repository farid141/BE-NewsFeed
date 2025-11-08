package router

import (
	"database/sql"

	"github.com/farid141/go-rest-api/controller"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, db *sql.DB) {
	setupPublicRoutes(app, db)

	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey:  jwtware.SigningKey{Key: []byte("secret")},
		TokenLookup: "cookie:token", // 👈 look in cookie named "token"
	}))

	setupAuthRoutes(app)
}

func setupAuthRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Get("/users", controller.GetUsers)
	api.Post("/users", controller.CreateUser)
}

func setupPublicRoutes(app *fiber.App, db *sql.DB) {
	api := app.Group("/api")

	api.Post("/login", controller.Login(db))
	api.Post("/logout", controller.Logout)
}
