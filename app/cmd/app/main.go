package main

import (
	"crypto/sha256"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/gofiber/template/html/v2"
	"github.com/mopeneko/blog-v2/app/model"
	"github.com/mopeneko/blog-v2/app/view"
	"github.com/mopeneko/blog-v2/app/view/dist"
	"github.com/mopeneko/blog-v2/app/view/tmpl"
)

func main() {
	engine := html.NewFileSystem(http.FS(tmpl.Content), ".html")

	loc, _ := time.LoadLocation("Asia/Tokyo")

	engine.AddFunc("date", func(t time.Time) string {
		return t.In(loc).Format("2006-01-02")
	})

	engine.AddFunc("unescape", func(s string) template.HTML {
		return template.HTML(s)
	})

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	css, err := dist.Content.ReadFile("style.css")
	if err != nil {
		log.Errorw("Failed to read CSS file", "err", err)
		return
	}

	cssHashBytes := sha256.Sum256(css)
	cssHash := fmt.Sprintf("%x", cssHashBytes)

	app.Use(logger.New())

	app.Use(func(c fiber.Ctx) error {
		c.Append("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		return c.Next()
	})

	app.Get("/dist/*", static.New("", static.Config{FS: dist.Content}))

	app.Get("/", func(c fiber.Ctx) error {
		articles, err := model.FetchArticles()
		if err != nil {
			log.Errorw("Failed to fetch articles", "err", err)
			return c.SendStatus(http.StatusInternalServerError)
		}

		return view.NewArticlesIndex(articles, cssHash).Render(c)
	})

	app.Get("/posts/:slug", func(c fiber.Ctx) error {
		article, err := model.FetchArticle(c.Params("slug"))
		if err != nil {
			log.Errorw("Failed to fetch article", "err", err)
			return c.SendStatus(http.StatusInternalServerError)
		}

		return view.NewArticle(article, cssHash).Render(c)
	})

	app.Get("/health", func(c fiber.Ctx) error {
		return c.SendStatus(http.StatusOK)
	})

	port := "3000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	log.Fatal(app.Listen(":" + port))
}
