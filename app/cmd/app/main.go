package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"html/template"
	"net/http"
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
	"github.com/mopeneko/blog-v2/app/public"
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

	engine.AddFunc("needToInsertAd", func(i int) bool {
		return (i+1)%3 == 0
	})

	app := fiber.New(fiber.Config{
		Views:         engine,
		ProxyHeader:   fiber.HeaderXForwardedFor,
		StrictRouting: true,
	})

	css, err := dist.Content.ReadFile("style.css")
	if err != nil {
		log.Errorw("Failed to read CSS file", "err", err)
		return
	}

	cssHashBytes := sha256.Sum256(css)
	cssHash := fmt.Sprintf("%x", cssHashBytes)

	client := model.NewArticleClient()
	pageClient := model.NewPageClient()

	app.Use(logger.New())

	app.Use(func(c fiber.Ctx) error {
		c.Append("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		return c.Next()
	})

	app.Get("/dist/*", static.New("", static.Config{FS: dist.Content}))

	app.Get("/public/*", static.New("", static.Config{FS: public.Content}))

	app.Get("/ads.txt", func(c fiber.Ctx) error {
		return c.SendString("google.com, pub-3857753364740983, DIRECT, f08c47fec0942fa0")
	})

	app.Get("/", func(c fiber.Ctx) error {
		articles, err := client.FetchArticles()
		if err != nil {
			log.Errorw("Failed to fetch articles", "err", err)
			return c.SendStatus(http.StatusInternalServerError)
		}

		return view.NewArticlesIndex(articles, cssHash).Render(c)
	})

	app.Get("/posts/:slug", func(c fiber.Ctx) error {
		article, err := client.FetchArticle(c.Params("slug"))
		if err != nil {
			log.Errorw("Failed to fetch article", "err", err)
			return c.SendStatus(http.StatusInternalServerError)
		}

		parsed, err := gohtml.Parse(strings.NewReader(article.Content))
		if err != nil {
			log.Errorw("Failed to parse HTML", "err", err)
			return c.SendStatus(http.StatusInternalServerError)
		}

		h2Count := 0

		var traverse func(n *gohtml.Node)
		traverse = func(n *gohtml.Node) {
			if n.Type == gohtml.ElementNode {
				if n.Data == "img" {
					n.Attr = append(n.Attr, gohtml.Attribute{
						Key: "loading",
						Val: "lazy",
					})
				}

				if n.Data == "h2" {
					h2Count++
					if h2Count == 2 {
						adNode := &gohtml.Node{
							Type: gohtml.ElementNode,
							Data: "div",
							Attr: []gohtml.Attribute{
								{Key: "class", Val: "ad-container"},
							},
						}

						// Create script node for AdSense JS
						scriptNode1 := &gohtml.Node{
							Type: gohtml.ElementNode,
							Data: "script",
							Attr: []gohtml.Attribute{
								{Key: "async", Val: ""},
								{Key: "src", Val: "https://pagead2.googlesyndication.com/pagead/js/adsbygoogle.js?client=ca-pub-3857753364740983"},
								{Key: "crossorigin", Val: "anonymous"},
							},
						}
						adNode.AppendChild(scriptNode1)

						// Create ins node for ad placement
						insNode := &gohtml.Node{
							Type: gohtml.ElementNode,
							Data: "ins",
							Attr: []gohtml.Attribute{
								{Key: "class", Val: "adsbygoogle"},
								{Key: "style", Val: "display:block; text-align:center;"},
								{Key: "data-ad-layout", Val: "in-article"},
								{Key: "data-ad-format", Val: "fluid"},
								{Key: "data-ad-client", Val: "ca-pub-3857753364740983"},
								{Key: "data-ad-slot", Val: "1281498636"},
							},
						}
						adNode.AppendChild(insNode)

						// Create script node for ad initialization
						scriptNode2 := &gohtml.Node{
							Type: gohtml.ElementNode,
							Data: "script",
						}
						scriptContent := &gohtml.Node{
							Type: gohtml.TextNode,
							Data: "(adsbygoogle = window.adsbygoogle || []).push({});",
						}
						scriptNode2.AppendChild(scriptContent)
						adNode.AppendChild(scriptNode2)

						if n.Parent != nil {
							n.Parent.InsertBefore(adNode, n)
						}
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

		relatedArticles, err := client.FetchArticlesByTags(article.Tags)
		if err != nil {
			log.Errorw("Failed to fetch related articles", "err", err)
			return c.SendStatus(http.StatusInternalServerError)
		}

		return view.NewArticle(article, cssHash, relatedArticles).Render(c)
	})

	app.Get("/posts/:slug/", func(c fiber.Ctx) error {
		return c.Redirect().Status(http.StatusMovedPermanently).To("/posts/" + c.Params("slug"))
	})

	app.Get("/pages/:slug", func(c fiber.Ctx) error {
		page, err := pageClient.FetchPage(c.Params("slug"))
		if err != nil {
			log.Errorw("Failed to fetch page", "err", err)
			return c.SendStatus(http.StatusInternalServerError)
		}

		return view.NewPage(page, cssHash).Render(c)
	})

	app.Get("/pages/:slug/", func(c fiber.Ctx) error {
		return c.Redirect().Status(http.StatusMovedPermanently).To("/pages/" + c.Params("slug"))
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
