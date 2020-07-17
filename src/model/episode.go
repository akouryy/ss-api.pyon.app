package model

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type Episode struct {
	Id        int       `json:"id"`
	Title     string    `json:"title"`
	BookId    int       `db:"book_id" json:"bookID"`
	IndexNum  int       `json:"index"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	Sections  []Section `json:"sections"`
}

func GetEpisode(dbx *sqlx.DB, episodeId int) (Episode, error) {
	var episode Episode
	err := dbx.Get(&episode, `SELECT * FROM episodes WHERE id = ?`, episodeId)
	if err != nil {
		return Episode{}, err
	}
	err = dbx.Select(&episode.Sections,
		`SELECT * FROM sections WHERE episode_id = ? ORDER BY indexnum`, episodeId)
	if err != nil {
		return Episode{}, err
	}
	return episode, nil
}

func GetEpisodes(dbx *sqlx.DB, bookId int) ([]Episode, error) {
	episodes := []Episode{}
	err := dbx.Select(&episodes,
		`SELECT * FROM episodes WHERE book_id = ? ORDER BY indexnum`, bookId)
	if err != nil {
		return nil, err
	}
	return episodes, nil
}
