package model

import (
	"github.com/mopeneko/blog-v2/app/newt"
)

type Article struct {
	newt.BaseContent

	Title     string     `json:"title"`
	Slug      string     `json:"slug"`
	Thumbnail newt.Image `json:"thumbnail"`
	Content   string     `json:"content"`
	Tags      []Tag      `json:"tags"`
}
