package newt

import "time"

type BaseContent struct {
	ID  string `json:"_id"`
	Sys struct {
		CreatedAt   time.Time `json:"createdAt"`
		UpdatedAt   time.Time `json:"updatedAt"`
		CustomOrder int       `json:"customOrder"`
		Raw         struct {
			CreatedAt        time.Time `json:"createdAt"`
			UpdatedAt        time.Time `json:"updatedAt"`
			FirstPublishedAt time.Time `json:"firstPublishedAt"`
			PublishedAt      time.Time `json:"publishedAt"`
		} `json:"raw"`
	} `json:"_sys"`
}

type Contents[T any] struct {
	Skip  int `json:"skip"`
	Limit int `json:"limit"`
	Total int `json:"total"`
	Items []T `json:"items"`
}

type Image struct {
	ID          string `json:"_id"`
	AltText     string `json:"altText"`
	Description string `json:"description"`
	FileName    string `json:"fileName"`
	FileSize    int    `json:"fileSize"`
	FileType    string `json:"fileType"`
	Height      int    `json:"height"`
	Src         string `json:"src"`
	Title       string `json:"title"`
	Width       int    `json:"width"`
}
