package model

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type Book struct {
	Id        int       `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}

func GetBooks(dbx *sqlx.DB) ([]Book, error) {
	books := []Book{}
	err := dbx.Select(&books, "SELECT * FROM books ORDER BY created_at DESC LIMIT 10;")
	return books, err
}
