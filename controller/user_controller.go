package controller

import (
	"strconv"

	"github.com/farid141/go-rest-api/service"
	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{userService}
}

func (ctl *UserController) GetUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	if limit == 0 {
		limit = 10
	}
	if page == 0 {
		page = 1
	}
	offset := (page - 1) * limit

	// claims := utils.GetJWTClaims(c)
	// userID := int(claims["id"].(float64))

	users, err := ctl.userService.ListUsers(1, page, limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(users)
}
