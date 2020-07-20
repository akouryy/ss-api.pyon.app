package model

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type Section struct {
	ID        int       `json:"id"`
	Content   string    `json:"content"`
	EpisodeID int       `db:"episode_id" json:"episodeID"`
	IndexNum  int       `json:"index"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

func CreateSection(dbx *sqlx.DB, content string, episodeID, indexNum int) error {
	_, err := dbx.Exec(
		`INSERT INTO sections(content, episode_id, indexnum, created_at) VALUES (?, ?, ?, NOW())`,
		content, episodeID, indexNum,
	)
	return err
}
