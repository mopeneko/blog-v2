package model

import "github.com/mopeneko/blog-v2/app/newt"

type Tag struct {
	newt.BaseContent
	Name string `json:"name"`
}
