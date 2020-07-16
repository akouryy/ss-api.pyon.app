package main

import (
	"fmt"
	"net/http"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

func booksHandler(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "{ message: \"Hello, world!\" }")
}

func main() {
	goji.Get("/books", booksHandler)
	goji.Serve()
}
