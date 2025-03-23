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
	bind := fiber.Map{
		"article": v.article,
		"cssHash": v.cssHash,
		"url":     "https://www.mope-blog.com/" + v.article.Slug,
		"title":   v.article.Title,
	}

	if v.article != nil && v.article.Thumbnail != nil {
		bind["image"] = v.article.Thumbnail.Src
	}

	return c.Render("article", bind, "layout")
}
