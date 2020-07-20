package model

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type Episode struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	BookID    int       `db:"book_id" json:"bookID"`
	IndexNum  int       `json:"index"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	Sections  []Section `json:"sections"`
}

// BookEpi ties an Episode with a Book.
type BookEpi struct {
	Book    Book    `json:"book"`
	Episode Episode `json:"episode"`
}

// GetEpisodeWithBook fetches an episode w/ its sections,
// along with the containing book w/o its episodes.
func GetEpisodeWithBook(dbx *sqlx.DB, episodeID int) (BookEpi, error) {
	var episode Episode
	err := dbx.Get(&episode, `SELECT * FROM episodes WHERE id = ?`, episodeID)
	if err != nil {
		return BookEpi{}, err
	}
	episode.Sections = []Section{}
	err = dbx.Select(&episode.Sections,
		`SELECT * FROM sections WHERE episode_id = ? ORDER BY indexnum`, episodeID)
	if err != nil {
		return BookEpi{}, err
	}
	book, err := GetBook(dbx, episode.BookID, false)
	if err != nil {
		return BookEpi{}, err
	}
	return BookEpi{book, episode}, nil
}

func GetEpisodes(dbx *sqlx.DB, bookID int) ([]Episode, error) {
	episodes := []Episode{}
	err := dbx.Select(&episodes,
		`SELECT * FROM episodes WHERE book_id = ? ORDER BY indexnum`, bookID)
	if err != nil {
		return nil, err
	}
	return episodes, nil
}

func CreateEpisode(dbx *sqlx.DB, title string, bookID, indexnum int) error {
	_, err := dbx.Exec(
		`INSERT INTO episodes(title, book_id, indexnum, created_at) VALUES (?, ?, ?, NOW())`,
		title, bookID, indexnum,
	)
	return err
}
