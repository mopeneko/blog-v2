package view

import (
	"slices"

	"github.com/gofiber/fiber/v3"
	"github.com/mopeneko/blog-v2/app/model"
)

type ArticleView struct {
	article         *model.Article
	cssHash         string
	relatedArticles []*model.Article
}

func NewArticle(article *model.Article, cssHash string, relatedArticles []*model.Article) *ArticleView {
	return &ArticleView{
		article:         article,
		cssHash:         cssHash,
		relatedArticles: relatedArticles,
	}
}

func (v *ArticleView) Render(c fiber.Ctx) error {
	relatedArticles := slices.DeleteFunc(v.relatedArticles, func(a *model.Article) bool {
		return a.Slug == v.article.Slug
	})

	if len(relatedArticles) > 3 {
		relatedArticles = relatedArticles[:3]
	}

	bind := fiber.Map{
		"article":         v.article,
		"cssHash":         v.cssHash,
		"url":             "https://www.mope-blog.com/posts/" + v.article.Slug,
		"title":           v.article.Title,
		"relatedArticles": relatedArticles,
	}

	if v.article != nil && v.article.Thumbnail != nil {
		bind["image"] = v.article.Thumbnail.Src
	}

	return c.Render("article", bind, "layout")
}
