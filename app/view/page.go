package view

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mopeneko/blog-v2/app/model"
)

type PageView struct {
	page    *model.Page
	cssHash string
}

func NewPage(page *model.Page, cssHash string) *PageView {
	return &PageView{
		page:    page,
		cssHash: cssHash,
	}
}

func (v *PageView) Render(c fiber.Ctx) error {
	bind := fiber.Map{
		"page":    v.page,
		"cssHash": v.cssHash,
		"url":     "https://www.mope-blog.com/pages/" + v.page.Slug,
		"title":   v.page.Title,
	}

	if v.page != nil && v.page.Thumbnail != nil {
		bind["image"] = v.page.Thumbnail.Src
	}

	return c.Render("page", bind, "layout")
}
