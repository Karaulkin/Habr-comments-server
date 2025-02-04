package models

import "time"

type Post struct {
	ID           int       `json:"id"`
	AuthorId     int       `json:"author_id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	AllowComment bool      `json:"allowComment"`
	Time         time.Time `json:"time"`
}
