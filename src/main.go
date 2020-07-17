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

	goji.Use(hutil.SetDBXMiddleware(dbx))
	goji.Post("/book", handler.BooksHandler)
	goji.Post("/author", handler.AuthorsHandler)
	goji.Post("/author/new", handler.NewAuthorHandler)
	goji.Serve()
}
