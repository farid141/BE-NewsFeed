package controller

import (
	"database/sql"

	"github.com/farid141/go-rest-api/helper"
	"github.com/farid141/go-rest-api/utils"
	"github.com/gofiber/fiber/v2"
)

func FollowUser(db *sql.DB, follow bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var err error
		followedID := c.Params("id")
		// followedID validation
		exists, err := helper.CoulmnValueExists(db, "users", "id", followedID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
			})
		}
		if !exists {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "ID not found",
			})
		}

		// insert new user

		claims := utils.GetJWTClaims(c)
		userID := claims["id"]
		queryString := func() string {
			if follow {
				return `INSERT INTO follows (follower_id, followed_id) VALUES (?, ?)`
			} else {
				return `DELETE FROM follows WHERE follower_id=? AND followed_id=?`
			}
		}()

		_, err = db.Exec(
			queryString,
			userID,
			followedID,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "Internal Server Error",
			})
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
