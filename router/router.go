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
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			if err != nil {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error": "access forbidden, invalid or expired token",
				})
			}
			return nil
		},
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
	api.Post("/register", controller.Register(db))
	api.Post("/refresh_token", controller.RefreshToken)
	api.Post("/logout", controller.Logout)
}
