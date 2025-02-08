package model

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	Location Location `json:"location"`
	Id uuid.UUID `json:"id"`
	Content string `json:"content"`
	Duration int64 `json:"duration"`
	CreatedAt int64 `json:"created_at"`
}

func NewPost(pr PostRequest) *Post {
	return &Post{
		Location: pr.Location,
		Id: uuid.New(),
		Content: pr.Content,
		Duration: pr.Duration,
		CreatedAt: time.Now().Unix(),
	}
}

func (p *Post) Valid() bool {
	return p.Location.Valid() && (len(p.Content) > 0) && ((p.Duration < 3600) && (p.Duration > 0))
}
