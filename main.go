package main

import (
	"github.com/farid141/go-rest-api/app"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	appContainer, err := app.InitializeApp()
	if err != nil {
		panic(err)
	}

	app := fiber.New(fiber.Config{
		// Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "Test App",
	})

	appContainer.Router.Setup(app)

	app.Use(cors.New(cors.Config{
		AllowOrigins:     appContainer.Config.ORIGINS,
		AllowHeaders:     "Origin, Content-Type, Accept",
		AllowCredentials: true,
	}))

	app.Listen(":3030")
}
