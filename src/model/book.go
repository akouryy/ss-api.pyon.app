package model

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type Book struct {
	Id        int       `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	Authors   []Author  `json:"authors"`
}

type bookAndAuthor struct {
	Book   Book
	Author Author
}

func GetBooks(dbx *sqlx.DB) ([]Book, error) {
	booksMap := map[int]*Book{}

	rows, err := dbx.Queryx(`
		SELECT
			books.id AS "book.id",
			books.title AS "book.title",
			books.created_at AS "book.created_at",
			authors.id AS "author.id",
			authors.name AS "author.name",
			authors.url AS "author.url",
			authors.created_at AS "author.created_at"
		FROM
			(SELECT * FROM books ORDER BY created_at DESC LIMIT 10) AS books
		INNER /*LEFT*/ JOIN
			book_authors ON books.id = book_authors.book_id
		INNER /*LEFT*/ JOIN
			authors ON authors.id = book_authors.author_id
	`)
	if err != nil {
		return nil, err
	}
	var row bookAndAuthor
	for rows.Next() {
		err := rows.StructScan(&row)
		if err != nil {
			return nil, err
		}
		if _, ok := booksMap[row.Book.Id]; !ok {
			booksMap[row.Book.Id] = &row.Book
		}
		book := booksMap[row.Book.Id]
		book.Authors = append(book.Authors, row.Author)
	}

	books := make([]Book, 0, len(booksMap))
	for _, book := range booksMap {
		books = append(books, *book)
	}
	return books, nil
}
