package routes

import (
	"encoding/json"
	"fmt"

	"github.com/NatanOih/urlShortenerGo/database"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

func ResolveUrl(c *fiber.Ctx) error {
	url := c.Params("url")

	r := database.CreateClient(0)
	defer r.Close()

	jsonData, err := r.Get(database.Ctx, url).Result()

	var data Data
	json.Unmarshal([]byte(jsonData), &data)

	value := data.URL
	clicks := data.Clicks

	fmt.Printf("URL: %s, Clicks: %d\n", value, clicks)

	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "short not found in the database"})
	} else if err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": " cannot connect to DB"})
	}

	data.Clicks += 1
	updatedJsonData, _ := json.Marshal(data)
	err = r.Set(database.Ctx, url, updatedJsonData, 0).Err() // 0 means keep the original TTL
	if err != nil {
		panic(err)
	}
	rInr := database.CreateClient(1)
	defer rInr.Close()

	_ = rInr.Incr(database.Ctx, "counter")

	return c.Redirect(value, 301)
}
