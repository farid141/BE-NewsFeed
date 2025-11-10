package controller

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/farid141/go-rest-api/dto"
	"github.com/farid141/go-rest-api/helper"
	"github.com/farid141/go-rest-api/utils"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

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
		var hashedPassword string
		err = db.QueryRow(
			`SELECT id, password_hash FROM users WHERE username=?`,
			req.Username,
		).Scan(&userID, &hashedPassword)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "invalid credentials"})
		}

		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password))
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
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

		var req dto.CreateUserRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		// username validation
		exists, err := helper.CoulmnValueExists(db, "users", "username", req.Username)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
			})
		}
		if exists {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"message": "Username already exists",
			})
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		// insert new user
		res, err := db.Exec(
			`INSERT INTO users (username, password_hash, created_at) VALUES (?,?,NOW())`,
			req.Username,
			hashedPassword,
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

		return c.Status(fiber.StatusCreated).JSON(dto.UserResponse{
			ID:       id,
			Username: req.Username,
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

func Me() fiber.Handler {
	return func(c *fiber.Ctx) error {

		claims := utils.GetJWTClaims(c)
		username := claims["username"]

		return c.JSON(fiber.Map{"name": username})
	}
}
