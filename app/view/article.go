package view

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mopeneko/blog-v2/app/model"
)

type ArticleView struct {
	article *model.Article
	cssHash string
}

func NewArticle(article *model.Article, cssHash string) *ArticleView {
	return &ArticleView{
		article: article,
		cssHash: cssHash,
	}
}

func (v *ArticleView) Render(c fiber.Ctx) error {
	return c.Render("article", fiber.Map{
		"article": v.article,
		"cssHash": v.cssHash,
	}, "layout")
}
