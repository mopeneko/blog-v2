package model

import "github.com/mopeneko/blog-v2/app/newt"

type Product struct {
	newt.BaseContent

	Name        string                    `json:"name"`
	Manufacture string                    `json:"manufacture"`
	Links       []newt.CustomField[*Link] `json:"links"`
	Image       *newt.Image               `json:"image"`
}
