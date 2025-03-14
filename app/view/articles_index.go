package view

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mopeneko/blog-v2/app/model"
)

type ArticlesIndexView struct {
	articles []*model.Article
}

func NewArticlesIndex(articles []*model.Article) *ArticlesIndexView {
	return &ArticlesIndexView{
		articles: articles,
	}
}

func (v *ArticlesIndexView) Render(c fiber.Ctx) error {
	return c.Render("articles_index", fiber.Map{
		"articles": v.articles,
	})
}
