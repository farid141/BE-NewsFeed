package controller

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/farid141/go-rest-api/utils"
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
			`SELECT id FROM users WHERE username=? AND password=?`,
			req.Username,
			req.Password,
		).Scan(&userID)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "invalid credentials"})
		}

		// token
		token, err := utils.GenerateJWT(req.Username, strconv.Itoa(userID), 15*time.Minute)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		// refresh_token
		c.Cookie(&fiber.Cookie{Name: "token", Value: token})
		refresh_token, err := utils.GenerateJWT(req.Username, strconv.Itoa(userID), 7*24*time.Hour)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		c.Cookie(&fiber.Cookie{Name: "refresh_token", Value: refresh_token})

		return c.JSON(fiber.Map{"token": token})
	}
}

func Register(db *sql.DB) fiber.Handler {
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

		res, err := db.Exec(
			`INSERT INTO users (username, password, createdat) VALUES (?,?,NOW())`,
			req.Username,
			req.Password,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
			})
		}

		id, err := res.LastInsertId()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
			})
		}

		return c.JSON(fiber.Map{
			"id":       id,
			"username": req.Username,
		})
	}
}

func RefreshToken(c *fiber.Ctx) error {
	refreshCookie := c.Cookies("refresh_token")
	if refreshCookie == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "no refresh token"})
	}

	token, err := jwt.Parse(refreshCookie, func(t *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid refresh token"})
	}

	claims := token.Claims.(jwt.MapClaims)
	username := claims["username"].(string)
	userId := claims["id"].(string)

	newAccessToken, err := utils.GenerateJWT(username, userId, 15*time.Minute)
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{Name: "token", Value: newAccessToken})

	return c.JSON(fiber.Map{"message": "token refreshed"})
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
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), // set expiry in the past
		HTTPOnly: true,
		Secure:   true, // set to true if using HTTPS
		SameSite: "Lax",
	})
	return c.SendStatus(fiber.StatusOK)
}
