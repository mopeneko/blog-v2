package newt

import "time"

type BaseContent struct {
	ID  string `json:"_id"`
	Sys struct {
		CreatedAt        time.Time `json:"createdAt"`
		UpdatedAt        time.Time `json:"updatedAt"`
		FirstPublishedAt time.Time `json:"firstPublishedAt"`
		PublishedAt      time.Time `json:"publishedAt"`
	} `json:"_sys"`
}

type Contents[T any] struct {
	Skip  int `json:"skip"`
	Limit int `json:"limit"`
	Total int `json:"total"`
	Items []T `json:"items"`
}
