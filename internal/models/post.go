package models

import "time"

type Post struct {
	ID            int
	AuthorId      int
	Title         string
	Content       string
	AllowComments bool
	Time          time.Time
}
