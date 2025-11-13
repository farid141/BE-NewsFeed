package controller

import (
	"github.com/farid141/go-rest-api/helper"
	"github.com/farid141/go-rest-api/service"
	"github.com/farid141/go-rest-api/utils"
	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{userService}
}

func (ctl *UserController) GetUsers(c *fiber.Ctx) error {
	page, limit := helper.GetPageAndLimit(c)
	offset := (page - 1) * limit

	claims := utils.GetJWTClaims(c)
	userID := claims["id"]
	id, ok := userID.(int)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}

	users, err := ctl.userService.ListUsers(id, page, limit, offset)
	if err != nil {
		if se, ok := err.(*helper.ServiceError); ok {
			return c.Status(se.StatusCode).JSON(fiber.Map{
				"error":   se.Message,
				"details": se.Details,
			})
		}
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}

	return c.JSON(users)
}

func (ctl *UserController) FollowUser(follow bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var err error
		followedID := c.Params("id")

		claims := utils.GetJWTClaims(c)
		userID := claims["id"]
		id, ok := userID.(string)
		if !ok {
			return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
		}

		err = ctl.userService.FollowUser(id, followedID, follow)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
		}

		return c.JSON(fiber.Map{
			"message": func() string {
				if follow {
					return "you are now following user " + followedID
				} else {
					return "you unfollowed user " + followedID
				}
			}(),
		})
	}
}
