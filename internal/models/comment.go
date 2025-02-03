package models

import "time"

type Comment struct {
	ID       int
	PostId   int
	AuthorId int
	ParentId int
	Content  []byte
	Time     time.Duration
}
