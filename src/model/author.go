package model

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type Author struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	URL       string    `json:"url"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

func GetAuthors(dbx *sqlx.DB) ([]Author, error) {
	authors := []Author{}
	err := dbx.Select(&authors, "SELECT * FROM authors ORDER BY created_at DESC LIMIT 10;")
	return authors, err
}

func CreateAuthor(dbx *sqlx.DB, name, url string) error {
	_, err := dbx.Exec(
		"INSERT INTO authors(name, url, created_at) VALUES (?, ?, NOW())",
		name, url,
	)
	return err
}
