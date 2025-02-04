package models

import "time"

type Comment struct {
	ID       int
	PostId   int
	AuthorId int
	ParentId *int
	Content  string
	Time     time.Time
}
