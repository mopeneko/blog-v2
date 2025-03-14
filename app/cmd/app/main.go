package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/template/html/v2"
	"github.com/mopeneko/blog-v2/app/model"
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
		articles := []*model.Article{
			{
				Title:       "Article 1",
				PublishedAt: time.Now(),
				UpdatedAt:   time.Now(),
				Tags:        []string{"tag1", "tag2"},
			},
			{
				Title:       "Article 2",
				PublishedAt: time.Now(),
				UpdatedAt:   time.Now(),
				Tags:        []string{"tag1", "tag2"},
			},
		}
		return view.NewArticlesIndex(articles).Render(c)
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
