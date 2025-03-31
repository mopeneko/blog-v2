package view

import (
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/mopeneko/blog-v2/app/model"
)

type ArticlesIndexView struct {
	articles []*model.Article
	cssHash  string
	isProd   bool
}

func NewArticlesIndex(articles []*model.Article, cssHash string) *ArticlesIndexView {
	return &ArticlesIndexView{
		articles: articles,
		cssHash:  cssHash,
		isProd:   os.Getenv("ENV") != "development",
	}
}

func (v *ArticlesIndexView) Render(c fiber.Ctx) error {
	return c.Render("articles_index", fiber.Map{
		"articles": v.articles,
		"cssHash":  v.cssHash,
		"url":      "https://www.mope-blog.com",
		"isProd":   v.isProd,
	}, "layout")
}
