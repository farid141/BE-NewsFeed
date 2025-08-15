package controller

import "github.com/gofiber/fiber/v2"

func GetUsers(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Get all users",
	})
}

func CreateUser(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "User created",
	})
}
