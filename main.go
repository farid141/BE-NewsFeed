package main

import (
	"fmt"
	"os"

	"github.com/farid141/go-rest-api/config"
	"github.com/farid141/go-rest-api/router"
	_ "github.com/go-sql-driver/mysql"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	dsn := fmt.Sprintf(
		"mysql://%s:%s@tcp(%s:%d)/%s",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)
	fmt.Println("DB URL: " + dsn)

	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "Test App",
	})

	file, _ := os.OpenFile("application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	log.SetOutput(file)

	router.SetupPublicRoutes(app)

	app.Get("/", func(c *fiber.Ctx) error {
		log.Info("hay")
		return c.SendString("Hello, World!")
	})

	// JWT Middleware
	app.Use(jwtware.New(jwtware.Config{
		SigningKey:  jwtware.SigningKey{Key: []byte("secret")},
		TokenLookup: "cookie:token", // 👈 look in cookie named "token"
	}))

	router.SetupAuthRoutes(app)

	app.Listen(":3000")
}
