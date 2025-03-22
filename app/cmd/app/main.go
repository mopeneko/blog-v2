package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	gohtml "golang.org/x/net/html"

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
		Views:       engine,
		ProxyHeader: "cf-connecting-ip",
	})

	css, err := dist.Content.ReadFile("style.css")
	if err != nil {
		log.Errorw("Failed to read CSS file", "err", err)
		return
	}

	cssHashBytes := sha256.Sum256(css)
	cssHash := fmt.Sprintf("%x", cssHashBytes)

	client := model.NewArticleClient()

	app.Use(logger.New())

	app.Use(func(c fiber.Ctx) error {
		c.Append("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		return c.Next()
	})

	app.Get("/dist/*", static.New("", static.Config{FS: dist.Content}))

	app.Get("/", func(c fiber.Ctx) error {
		articles, err := client.FetchArticles()
		if err != nil {
			log.Errorw("Failed to fetch articles", "err", err)
			return c.SendStatus(http.StatusInternalServerError)
		}

		for i := range articles {
			if articles[i].Thumbnail == nil {
				continue
			}

			thumbSrc := articles[i].Thumbnail.Src
			if thumbSrc == "" {
				continue
			}

			u, err := url.Parse(thumbSrc)
			if err != nil {
				log.Errorw("Failed to parse thumbnail URL", "err", err)
				continue
			}

			q := u.Query()
			q.Add("format", "auto")
			u.RawQuery = q.Encode()
			articles[i].Thumbnail.Src = u.String()
		}

		return view.NewArticlesIndex(articles, cssHash).Render(c)
	})

	app.Get("/posts/:slug", func(c fiber.Ctx) error {
		if strings.HasSuffix(c.OriginalURL(), "/") {
			return c.Redirect().Status(http.StatusMovedPermanently).To("/posts/" + c.Params("slug"))
		}

		article, err := client.FetchArticle(c.Params("slug"))
		if err != nil {
			log.Errorw("Failed to fetch article", "err", err)
			return c.SendStatus(http.StatusInternalServerError)
		}

		if article.Thumbnail != nil && article.Thumbnail.Src != "" {
			u, err := url.Parse(article.Thumbnail.Src)
			if err != nil {
				log.Errorw("Failed to parse thumbnail URL", "err", err)
				return c.SendStatus(http.StatusInternalServerError)
			}

			q := u.Query()
			q.Add("format", "auto")
			u.RawQuery = q.Encode()
			article.Thumbnail.Src = u.String()
		}

		parsed, err := gohtml.Parse(strings.NewReader(article.Content))
		if err != nil {
			log.Errorw("Failed to parse HTML", "err", err)
			return c.SendStatus(http.StatusInternalServerError)
		}

		var traverse func(n *gohtml.Node)
		traverse = func(n *gohtml.Node) {
			if n.Type == gohtml.ElementNode && n.Data == "img" {
				for i := range n.Attr {
					if n.Attr[i].Key == "src" {
						u, err := url.Parse(n.Attr[i].Val)
						if err != nil {
							log.Errorw("Failed to parse image URL", "err", err)
							continue
						}

						q := u.Query()
						q.Add("format", "auto")
						u.RawQuery = q.Encode()
						n.Attr[i].Val = u.String()
					}
				}
			}

			for c := n.FirstChild; c != nil; c = c.NextSibling {
				traverse(c)
			}
		}

		traverse(parsed)

		var buf bytes.Buffer
		gohtml.Render(&buf, parsed)
		article.Content = buf.String()

		if article.Product != nil && article.Product.Image != nil && article.Product.Image.Src != "" {
			u, err := url.Parse(article.Product.Image.Src)
			if err != nil {
				log.Errorw("Failed to parse thumbnail URL", "err", err)
				return c.SendStatus(http.StatusInternalServerError)
			}

			q := u.Query()
			q.Add("format", "auto")
			u.RawQuery = q.Encode()
			article.Product.Image.Src = u.String()
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
