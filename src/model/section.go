package model

import (
	"time"
)

type Section struct {
	Id        int       `json:"id"`
	Content   string    `json:"content"`
	EpisodeID int       `db:"episode_id" json:"episodeID"`
	IndexNum  int       `json:"index"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}
