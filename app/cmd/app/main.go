package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/template/html/v2"
	"github.com/mopeneko/blog-v2/app/model"
	"github.com/mopeneko/blog-v2/app/newt"
	"github.com/mopeneko/blog-v2/app/view"
	"github.com/mopeneko/blog-v2/app/view/tmpl"
)

func main() {
	engine := html.NewFileSystem(http.FS(tmpl.Content), ".html")

	loc, _ := time.LoadLocation("Asia/Tokyo")

	engine.AddFunc("date", func(t time.Time) string {
		return t.In(loc).Format("2006-01-02")
	})

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", func(c fiber.Ctx) error {
		req, err := http.NewRequest("GET", os.Getenv("NEWT_API"), nil)
		if err != nil {
			log.Errorw("Failed to create request", "err", err)
			return c.SendStatus(http.StatusInternalServerError)
		}

		req.Header.Set("Authorization", "Bearer "+os.Getenv("NEWT_TOKEN"))

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Errorw("Failed to send request", "err", err)
			return c.SendStatus(http.StatusInternalServerError)
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Errorw("Failed to get articles", "status", resp.StatusCode)
			return c.SendStatus(http.StatusInternalServerError)
		}

		contents := new(newt.Contents[*model.Article])
		if err := json.NewDecoder(resp.Body).Decode(contents); err != nil {
			log.Errorw("Failed to decode response", "err", err)
			return c.SendStatus(http.StatusInternalServerError)
		}
		return view.NewArticlesIndex(contents.Items).Render(c)
	})

	host := "localhost"
	if os.Getenv("HOST") != "" {
		host = os.Getenv("HOST")
	}

	port := "3000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	log.Fatal(app.Listen(host + ":" + port))
}
