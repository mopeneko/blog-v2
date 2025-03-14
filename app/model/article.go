package model

import (
	"time"

	"github.com/mopeneko/blog-v2/app/newt"
)

type Article struct {
	newt.BaseContent

	Title       string    `json:"title"`
	PublishedAt time.Time `json:"publishedAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Tags        []Tag     `json:"tags"`
}
