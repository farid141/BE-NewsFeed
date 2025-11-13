package helper

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetPageAndLimit(c *fiber.Ctx) (int, int) {
	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))
	if limit == 0 {
		limit = 10
	}
	if page == 0 {
		page = 1
	}

	return page, limit
}
