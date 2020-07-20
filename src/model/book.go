package model

import (
	"errors"
	"sort"
	"time"

	"github.com/akouryy/ss-api.pyon.app/src/model/mutil"
	"github.com/jmoiron/sqlx"
)

// Book is a nutritious book.
type Book struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	Authors   []Author  `json:"authors"`
	Episodes  []Episode `json:"episodes"`
}

type bookAndAuthor struct {
	Book   Book
	Author Author
}

// GetBook fetches an book with its authors and optionally with its episodes.
func GetBook(dbx *sqlx.DB, bookID int, withEpisode bool) (Book, error) {
	var book Book
	err := dbx.Get(&book, `SELECT * FROM books WHERE id = ?`, bookID)
	if err != nil {
		return Book{}, err
	}
	err = dbx.Select(&book.Authors, `
		SELECT authors.* FROM authors
		INNER JOIN
			(SELECT * FROM book_authors WHERE book_id = ?) AS book_authors
			ON authors.id = book_authors.author_id
	`, bookID)
	if err != nil {
		return Book{}, err
	}
	if withEpisode {
		book.Episodes, err = GetEpisodes(dbx, bookID)
		if err != nil {
			return Book{}, err
		}
	}
	return book, nil
}

// GetBooks fetches the latest books with their authors but without their episodes.
func GetBooks(dbx *sqlx.DB) ([]Book, error) {
	booksMap := map[int]*Book{}

	rows, err := dbx.Queryx(`
		SELECT
			books.id AS "book.id", books.title AS "book.title", books.created_at AS "book.created_at",
			authors.id AS "author.id", authors.name AS "author.name", 
				authors.url AS "author.url", authors.created_at AS "author.created_at"
		FROM
			(SELECT * FROM books ORDER BY created_at DESC LIMIT 10) AS books
		INNER /*LEFT*/ JOIN book_authors
			ON books.id = book_authors.book_id
		INNER /*LEFT*/ JOIN authors
			ON authors.id = book_authors.author_id
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
		if _, ok := booksMap[row.Book.ID]; !ok {
			booksMap[row.Book.ID] = &row.Book
		}
		book := booksMap[row.Book.ID]
		book.Authors = append(book.Authors, row.Author)
	}

	books := make([]Book, 0, len(booksMap))
	for _, book := range booksMap {
		books = append(books, *book)
	}
	sort.Slice(books, func(i, j int) bool {
		return books[i].CreatedAt. /*IsAfter*/ After(books[j].CreatedAt)
	})
	return books, nil
}

// AddBookAuthor registers an existing Author as an author of an existing Book.
func AddBookAuthor(dbx *sqlx.DB, bookID, authorID int) error {
	_, err := dbx.Exec(
		`INSERT INTO book_authors(book_id, author_id) VALUES (?, ?)`,
		bookID, authorID,
	)
	return err
}

// RemoveBookAuthor unregisters an Author from the authors of a Book.
func RemoveBookAuthor(dbx *sqlx.DB, bookID, authorID int) error {
	return mutil.Transaction(dbx, func(tx *sqlx.Tx) error {
		_, err := dbx.Exec(`SELECT id FROM books WHERE id = ? FOR UPDATE`, bookID)
		if err != nil {
			return err
		}

		var cnt int
		err = dbx.Get(&cnt, `SELECT COUNT(*) FROM book_authors WHERE book_id = ?`, bookID)
		if err != nil {
			return err
		}
		if cnt <= 1 {
			return errors.New("there must remain some authors for this book")
		}

		_, err = dbx.Exec(
			`DELETE FROM book_authors WHERE book_id = ? AND author_id = ?`, bookID, authorID,
		)

		return err
	})
}
