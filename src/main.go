package main

import (
	"log"
	"os"

	"github.com/akouryy/ss-api.pyon.app/src/handler"
	"github.com/akouryy/ss-api.pyon.app/src/handler/hutil"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/zenazn/goji"
)

var dbx *sqlx.DB

func main() {
	var err error

	dbx, err = sqlx.Connect("mysql", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalln(err)
	}

	// goji.Use(hutil.SetDBXMiddleware(dbx))
	goji.Post("/author/list", hutil.Wrap(dbx, handler.AuthorsHandler))
	goji.Post("/author/new", hutil.Wrap(dbx, handler.NewAuthorHandler))
	goji.Post("/book/list", hutil.Wrap(dbx, handler.BooksHandler))
	goji.Post("/book/author/new", hutil.Wrap(dbx, handler.NewBookAuthorHandler))
	goji.Post("/book/author/delete", hutil.Wrap(dbx, handler.DeleteBookAuthorHandler))
	goji.Serve()
}
