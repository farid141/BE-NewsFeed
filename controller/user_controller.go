package controller

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/farid141/go-rest-api/model"
	"github.com/farid141/go-rest-api/response"
	"github.com/gofiber/fiber/v2"
)

func GetUsers(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
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
            SELECT id, username, password, createdat
			FROM users
            ORDER BY createdat DESC
            LIMIT ? OFFSET ?
        `, limit, offset)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		defer rows.Close()

		var users []model.User
		for rows.Next() {
			var u model.User
			var createdAtStr string

			if err := rows.Scan(&u.ID, &u.Username, &u.Password, &createdAtStr); err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}

			u.CreatedAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
			if err != nil {
				return c.Status(500).JSON(fiber.Map{"error": err.Error()})
			}
			users = append(users, u)
		}

		resp := response.PaginatedResponse[model.User]{
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
