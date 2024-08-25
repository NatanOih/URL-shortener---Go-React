package main

import (
	"fmt"
	"log"

	"github.com/NatanOih/urlShortenerGo/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func setupRoutes(app *fiber.App) {
	app.Get("/:url", routes.ResolveUrl)
	app.Get("api/v1/url-list", routes.GetAllUrls)
	app.Post("/api/v1", routes.ShortenURL)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}
	app := fiber.New()

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET, POST, PUT, DELETE",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	fmt.Print("app is running on port 3000")

	app.Static("/", "./frontend/dist")

	setupRoutes(app)

	// log.Fatal(app.Listen(os.Getenv("APP_PORT")))
	log.Fatal(app.Listen(":3000"))

}
