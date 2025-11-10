package controller

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/farid141/go-rest-api/dto"
	"github.com/farid141/go-rest-api/model"
	"github.com/farid141/go-rest-api/response"
	"github.com/farid141/go-rest-api/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func CreatePost(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var validate = validator.New()
		var err error

		var req dto.CreatePostRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": err.Error()})
		}

		// Validate struct
		if err := validate.Struct(req); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		claims := utils.GetJWTClaims(c)
		userID := claims["id"]

		res, err := db.Exec(
			`INSERT INTO posts (user_id, content, created_at) VALUES (?,?,NOW())`,
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
		err = db.QueryRow("SELECT id, user_id, content, created_at FROM posts WHERE id = ?", lastID).
			Scan(&createdPost.ID, &createdPost.UserID, &createdPost.Content, &createdAtStr)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": "invalid credentials"})
		}

		createdPost.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusCreated).JSON(dto.PostResponse{
			ID:        lastID,
			UserID:    createdPost.UserID,
			Content:   createdPost.Content,
			CreatedAt: createdPost.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
}

func GetFeed(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims := utils.GetJWTClaims(c)
		userID := claims["id"]

		page, _ := strconv.Atoi(c.Query("page"))
		limit, _ := strconv.Atoi(c.Query("limit"))
		if limit == 0 {
			limit = 10
		}
		if page == 0 {
			page = 1
		}
		offset := (page - 1) * limit

		// ambil total dulu untuk pagination
		var total int
		err := db.QueryRow(`
            SELECT COUNT(*)
            FROM posts p
            LEFT JOIN follows f ON f.followed_id = p.user_id
            WHERE f.follower_id = ? OR p.user_id = ?
        `, userID, userID).Scan(&total)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		rows, err := db.Query(`
            SELECT p.id, p.user_id, p.content, p.created_at
            FROM posts p
            LEFT JOIN follows f ON f.followed_id = p.user_id
            WHERE f.follower_id = ? OR p.user_id = ?
            ORDER BY p.created_at DESC
            LIMIT ? OFFSET ?
        `, userID, userID, limit, offset)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		defer rows.Close()

		posts := make([]dto.PostResponse, 0)
		for rows.Next() {
			var p dto.PostResponse
			if err := rows.Scan(&p.ID, &p.UserID, &p.Content, &p.CreatedAt); err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}

			posts = append(posts, p)
		}

		resp := response.PaginatedResponse[dto.PostResponse]{
			Data: posts,
			Pagination: response.Pagination{
				Page:    page,
				Limit:   limit,
				Total:   total,
				HasMore: page*limit < total,
			},
		}

		return c.JSON(resp)
	}
}
