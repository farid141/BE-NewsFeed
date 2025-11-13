package controller

import (
	"time"

	"github.com/farid141/go-rest-api/dto"
	"github.com/farid141/go-rest-api/helper"
	"github.com/farid141/go-rest-api/service"
	"github.com/farid141/go-rest-api/utils"
	"github.com/gofiber/fiber/v2"

	"github.com/golang-jwt/jwt/v5"
)

type AuthController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{authService}
}

func (ctl *AuthController) Login(c *fiber.Ctx) error {
	var err error

	var req dto.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	token, err := ctl.authService.Login(req.Username, req.Password)
	if err != nil {
		if se, ok := err.(*helper.ServiceError); ok {
			return c.Status(se.StatusCode).JSON(fiber.Map{
				"error":   se.Message,
				"details": se.Details,
			})
		}
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}

	c.Cookie(&fiber.Cookie{Name: "token", Value: token.AccessToken})
	c.Cookie(&fiber.Cookie{Name: "token", Value: token.RefreshToken})
	return c.JSON(fiber.Map{"token": token.AccessToken})
}

func (ctl *AuthController) Register(c *fiber.Ctx) error {
	var err error

	var req dto.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	var user *dto.CreateUserResponse
	user, err = ctl.authService.Register(req)

	if err != nil {
		if se, ok := err.(*helper.ServiceError); ok {
			return c.Status(se.StatusCode).JSON(fiber.Map{
				"error":   se.Message,
				"details": se.Details,
			})
		}
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}
	return c.Status(fiber.StatusCreated).JSON(user)
}

func (ctl *AuthController) RefreshToken(c *fiber.Ctx) error {
	refreshCookie := c.Cookies("refresh_token")
	if refreshCookie == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "no refresh token"})
	}

	// Input validation
	token, err := jwt.Parse(refreshCookie, func(t *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid refresh token"})
	}

	// Process
	claims := token.Claims.(jwt.MapClaims)
	username := claims["username"].(string)
	userId := claims["id"].(string)
	newAccessToken, err := utils.GenerateJWT(username, userId, 15*time.Minute)
	if err != nil {
		return err
	}

	// structuring response
	c.Cookie(&fiber.Cookie{Name: "token", Value: newAccessToken})

	return c.JSON(fiber.Map{"message": "token refreshed"})
}

func (ctl *AuthController) Logout(c *fiber.Ctx) error {
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

func (ctl *AuthController) Me(c *fiber.Ctx) error {
	claims := utils.GetJWTClaims(c)
	username := claims["username"]

	return c.JSON(fiber.Map{"name": username})
}
