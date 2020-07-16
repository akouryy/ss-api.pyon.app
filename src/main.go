package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

var db *sqlx.DB

type Book struct {
	Id        string
	Title     string
	CreatedAt time.Time `db:"created_at"`
}

func booksHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	books := []Book{}
	err := db.Select(&books, "SELECT * FROM books ORDER BY created_at DESC LIMIT 10;")
	if err != nil {
		fmt.Fprintf(w, "error %s", err.Error())
		return
	}

	j, err := json.Marshal(books)
	if err != nil {
		fmt.Fprintf(w, "error %s", err.Error())
		return
	}
	io.WriteString(w, string(j))
}

func main() {
	var err error

	db, err = sqlx.Connect("mysql", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalln(err)
	}

	goji.Get("/books", booksHandler)
	goji.Serve()
}
