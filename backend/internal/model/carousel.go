package model

import "time"

type CarouselImage struct {
	ID        int64     `json:"id"`
	ImageURL  string    `json:"image_url"`
	Title     string    `json:"title"`
	LinkURL   string    `json:"link_url"`
	SortOrder int       `json:"sort_order"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
