package main

import (
	"i9pxc/globals"
	"i9pxc/initializers"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
)

func init() {
	err := initializers.Init()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	app := fiber.New()

	app.Use(cache.New(cache.Config{
		Storage:              globals.RedisStore,
		Expiration:           2 * time.Minute,
		StoreResponseHeaders: true,
	}))

	app.Get("*", func(c *fiber.Ctx) error {
		url := c.OriginalURL()

		// agent := fiber.Get(os.Getenv("ORIGIN_SERVER") + url)

		res, err := http.Get(os.Getenv("ORIGIN_SERVER") + url)
		if err != nil {
			log.Println(err)
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		for hKey, hVals := range res.Header {
			for i, hVal := range hVals {
				if i == 0 {
					c.Set(hKey, hVal)
					continue
				}
				c.Append(hKey, hVal)
			}
		}

		defer res.Body.Close()

		data, err := io.ReadAll(res.Body)
		if err != nil {
			log.Println(err)
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return c.Status(res.StatusCode).Send(data)
	})

	log.Fatalln(app.Listen("localhost:" + "5000"))
}
