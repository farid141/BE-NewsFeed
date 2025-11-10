package controller

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/farid141/go-rest-api/dto"
	"github.com/farid141/go-rest-api/response"
	"github.com/farid141/go-rest-api/utils"
	"github.com/gofiber/fiber/v2"
)

func GetUsers(db *sql.DB) fiber.Handler {
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
            FROM users
        `).Scan(&total)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		rows, err := db.Query(`
            SELECT 
				u.id,
				u.username,
				CASE WHEN f.follower_id IS NOT NULL THEN TRUE ELSE FALSE END AS is_following,
				u.created_at
			FROM users u
			LEFT JOIN follows f ON f.followed_id = u.id AND f.follower_id = ?
			WHERE u.id != ?
			ORDER BY u.created_at DESC
            LIMIT ? OFFSET ?
        `, userID, userID, limit, offset)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		defer rows.Close()

		users := make([]dto.UserResponse, 0)
		for rows.Next() {
			var u dto.UserResponse
			var createdAtStr string

			if err := rows.Scan(&u.ID, &u.Username, &u.Following, &createdAtStr); err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}

			u.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}
			users = append(users, u)
		}

		resp := response.PaginatedResponse[dto.UserResponse]{
			Data: users,
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
