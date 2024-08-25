package routes

import (
	"github.com/NatanOih/urlShortenerGo/database"
	"github.com/NatanOih/urlShortenerGo/helpers"
	"github.com/gofiber/fiber/v2"
)

func GetAllUrls(c *fiber.Ctx) error {
	r := database.CreateClient(0)
	defer r.Close()

	ctx := database.Ctx

	JsonForUi, err := helpers.FetchDataFromRedis(ctx, r)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to fetch data from Redis")
	}

	return c.JSON(JsonForUi)
}
