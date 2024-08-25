package routes

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/NatanOih/urlShortenerGo/database"
	"github.com/NatanOih/urlShortenerGo/helpers"
	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"CustomShort"`
	Expiry      time.Duration `json:"expiry"`
}

type response struct {
	URL             string        `json:"url"`
	CustomShort     string        `json:"short"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"rate_limit"`
	XRateLimitReset time.Duration `json:"rate_limit_reset"`
	Clicks          int           `json:"clicks"`
}

type Data struct {
	URL    string `json:"url"`
	Clicks int    `json:"clicks"`
}

func ShortenURL(c *fiber.Ctx) error {
	body := new(request)
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
	}
	fmt.Printf("Got a new request: %+v\n", body)

	bodyBytes := c.Body()
	fmt.Printf("Body: %s\n", bodyBytes)

	//implement rate limiting, check the ip of the caller and check if the user used our rate.

	r2 := database.CreateClient(1)
	defer r2.Close()
	val, err := r2.Get(database.Ctx, c.IP()).Result()
	if err == redis.Nil {
		_ = r2.Set(database.Ctx, c.IP(), os.Getenv("API_QUOTA"), 30*60*time.Second)
	} else {

		valInt, _ := strconv.Atoi(val)
		if valInt <= 0 {
			limit, _ := r2.TTL(database.Ctx, c.IP()).Result()
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error":                    "Rate limit exceeded",
				"rate_limit_time_to_reset": limit / time.Nanosecond / time.Minute,
			})
		}
	}

	//check if the input is an actual url

	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid URL"})
	}

	// check for domain error

	if !helpers.RemoveDomainError(body.URL) {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "domain loop"})
	}

	//enforce https, ssl

	body.URL = helpers.EnforceHTTP(body.URL)

	var id string

	fmt.Print("custom short is \n", body.CustomShort)

	if body.CustomShort == "" {
		id = uuid.New().String()[:6]
	} else {
		id = body.CustomShort
	}

	r := database.CreateClient(0)
	defer r.Close()

	val, _ = r.Get(database.Ctx, id).Result()
	if val != "" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "URL custom short is already in use dick natan",
		})
	}

	if body.Expiry == 0 {
		body.Expiry = 24
	}

	data := Data{
		URL:    body.URL,
		Clicks: 0,
	}

	jsonData, _ := json.Marshal(data)

	err = r.Set(database.Ctx, id, jsonData, body.Expiry*3600*time.Second).Err()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": " Unable to connect to serverccc33",
		})
	}

	resp := response{
		URL:             data.URL,
		Clicks:          data.Clicks,
		CustomShort:     "",
		Expiry:          body.Expiry,
		XRateRemaining:  10,
		XRateLimitReset: 30,
	}

	r2.Decr(database.Ctx, c.IP())

	val, _ = r2.Get(database.Ctx, c.IP()).Result()
	resp.XRateRemaining, _ = strconv.Atoi(val)

	ttl, _ := r2.TTL(database.Ctx, c.IP()).Result()

	resp.XRateLimitReset = ttl / time.Nanosecond / time.Minute

	resp.CustomShort = os.Getenv("DOMAIN") + "/" + id

	return c.Status(fiber.StatusOK).JSON(resp)

}
