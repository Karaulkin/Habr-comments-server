package models

import "time"

type Comment struct {
	ID       int       `json:"id"`
	PostId   int       `json:"post_id"`
	AuthorId int       `json:"author_id"`
	ParentId int       `json:"parent_id"`
	Content  string    `json:"content"`
	Time     time.Time `json:"time"`
}
