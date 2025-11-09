package controller

import (
	"database/sql"
	"time"

	"github.com/farid141/go-rest-api/dto"
	"github.com/farid141/go-rest-api/model"
	"github.com/farid141/go-rest-api/utils"
	"github.com/gofiber/fiber/v2"
)

func CreatePost(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var err error

		var req dto.CreatePostRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		claims := utils.GetJWTClaims(c)
		userID := claims["id"]

		res, err := db.Exec(
			`INSERT INTO posts (userid, content, createdat) VALUES (?,?,NOW())`,
			userID,
			req.Content,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
			})
		}

		lastID, err := res.LastInsertId()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
			})
		}

		var createdAtStr string
		var createdPost model.Post
		err = db.QueryRow("SELECT id, userid, content, createdat FROM posts WHERE id = ?", lastID).
			Scan(&createdPost.ID, &createdPost.UserID, &createdPost.Content, &createdAtStr)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "invalid credentials"})
		}

		createdPost.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(dto.PostResponse{
			ID:        lastID,
			UserID:    createdPost.UserID,
			Content:   createdPost.Content,
			CreatedAt: createdPost.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
}
