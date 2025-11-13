package controller

import (
	"github.com/farid141/go-rest-api/dto"
	"github.com/farid141/go-rest-api/helper"
	"github.com/farid141/go-rest-api/service"
	"github.com/farid141/go-rest-api/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type PostController struct {
	postService service.PostService
}

func NewPostController(postService service.PostService) *PostController {
	return &PostController{postService}
}

func (ctl *PostController) CreatePost(c *fiber.Ctx) error {
	var req dto.CreatePostRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{"error": err.Error()})
	}

	// Validate struct
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	claims := utils.GetJWTClaims(c)
	userID := claims["id"]
	id, ok := userID.(string)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}

	post, err := ctl.postService.CreatePost(id, req)
	if err != nil {
		if se, ok := err.(*helper.ServiceError); ok {
			return c.Status(se.StatusCode).JSON(fiber.Map{
				"error":   se.Message,
				"details": se.Details,
			})
		}
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}

	return c.JSON(post)
}

func (ctl *PostController) GetFeed(c *fiber.Ctx) error {
	page, limit := helper.GetPageAndLimit(c)
	offset := (page - 1) * limit

	claims := utils.GetJWTClaims(c)
	userID := claims["id"]
	id, ok := userID.(string)
	if !ok {
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}

	posts, err := ctl.postService.GetFeed(id, page, limit, offset)
	if err != nil {
		if se, ok := err.(*helper.ServiceError); ok {
			return c.Status(se.StatusCode).JSON(fiber.Map{
				"error":   se.Message,
				"details": se.Details,
			})
		}
		return c.Status(500).JSON(fiber.Map{"error": "internal server error"})
	}

	return c.JSON(posts)
}
