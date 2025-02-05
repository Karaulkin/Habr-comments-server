package models

import "time"

type Post struct {
	ID            int       `json:"id"`
	AuthorId      int       `json:"authorId"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	AllowComments bool      `json:"allowComments"`
	CreatedAt     time.Time `json:"createdAt"`
}
