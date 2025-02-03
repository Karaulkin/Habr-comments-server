package models

import "time"

type Post struct {
	ID           int
	AuthorId     int
	Title        string
	Content      []byte
	AllowComment bool
	Time         time.Duration
}
