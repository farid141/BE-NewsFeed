package controller

import (
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Login(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var err error

		type Req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		var req Req
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		var userID int
		err = db.QueryRow(
			`SELECT id FROM users WHERE name=? AND password=?`,
			req.Username,
			req.Password,
		).Scan(&userID)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "invalid credentials"})
		}

		// Create the Claims
		claims := jwt.MapClaims{
			"name": req.Username,
			"id":   userID,
			"exp":  time.Now().Add(time.Hour * 72).Unix(),
		}

		// Create token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		c.Cookie(&fiber.Cookie{Name: "token", Value: t})
		return c.JSON(fiber.Map{"token": t})
	}
}

func Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), // set expiry in the past
		HTTPOnly: true,
		Secure:   true, // set to true if using HTTPS
		SameSite: "Lax",
	})
	return c.SendStatus(fiber.StatusOK)
}
