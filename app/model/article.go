package model

import "time"

type Article struct {
	Title       string
	PublishedAt time.Time
	UpdatedAt   time.Time
	Tags        []string
}
