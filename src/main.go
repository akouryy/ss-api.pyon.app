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
	"golang.org/x/crypto/bcrypt"
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

func BooksHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user, err := Authenticate(r)
	if ReportError(w, err) {
		return
	}
	log.Println(user.Nickname)

	books := []Book{}
	if ReportError(w,
		dbx.Select(&books, "SELECT * FROM books ORDER BY created_at DESC LIMIT 10;"),
	) {
		return
	}

	j, err := json.Marshal(books)
	if ReportError(w, err) {
		return
	}
	io.WriteString(w, string(j))
}

func main() {
	var err error

	dbx, err = sqlx.Connect("mysql", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalln(err)
	}

	goji.Post("/books", BooksHandler)
	goji.Serve()

	a, err := bcrypt.GenerateFromPassword([]byte("OYXgulIdzvqcrT7rusVGJtMEHspIFDIwB6qltLSEXjiS4wMnEUDcVLMi6FXLRC2Yo6DKej"), 4)
	fmt.Println(string(a), err)
}
