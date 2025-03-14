package view

import (
	"github.com/gofiber/fiber/v3"
)

type ArticlesIndexView struct {
}

func NewArticlesIndex() *ArticlesIndexView {
	return &ArticlesIndexView{}
}

func (v *ArticlesIndexView) Render(c fiber.Ctx) error {
	return c.Render("articles_index", fiber.Map{})
}
