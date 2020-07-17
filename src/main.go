package main

import (
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

var dbx *sqlx.DB

type User struct {
	Id       int
	Nickname string
	Password []byte
}

type Book struct {
	Id        int
	Title     string
	CreatedAt time.Time `db:"created_at"`
}

func BooksHandler(c web.C, w http.ResponseWriter, httpReq *http.Request) {
	body, err := ReadWholeBody(httpReq)
	if ReportError(w, err) {
		return
	}

	_, err = Authenticate(body)
	if ReportError(w, err) {
		return
	}

	books := []Book{}
	if ReportError(w,
		dbx.Select(&books, "SELECT * FROM books ORDER BY created_at DESC LIMIT 10;"),
	) {
		return
	}

	RenderJSON(w, books)
}

func main() {
	var err error

	dbx, err = sqlx.Connect("mysql", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalln(err)
	}

	goji.Post("/book", BooksHandler)
	goji.Post("/author", AuthorsHandler)
	goji.Post("/author/new", NewAuthorHandler)
	goji.Serve()
}
