package handler

import (
	"encoding/json"
	"net/http"

	"github.com/akouryy/ss-api.pyon.app/src/handler/hutil"
	"github.com/akouryy/ss-api.pyon.app/src/model"
	"github.com/jmoiron/sqlx"
	"github.com/zenazn/goji/web"
)

type reqBook struct {
	BookID int `json:"bookID"`
}

// BookHandler handles /book/show and responds with a book.
func BookHandler(
	ctx web.C, writer http.ResponseWriter, httpReq *http.Request,
	dbx *sqlx.DB, body []byte, _ model.User,
) {
	var req reqBook
	err := json.Unmarshal(body, &req)
	if hutil.ReportError(writer, err) {
		return
	}

	book, err := model.GetBook(dbx, req.BookID, true)
	if hutil.ReportError(writer, err) {
		return
	}

	hutil.RenderJSON(writer, book)
}

// BooksHandler handles /book/list and responds with the latest books.
func BooksHandler(
	ctx web.C, writer http.ResponseWriter, httpReq *http.Request,
	dbx *sqlx.DB, body []byte, _ model.User,
) {
	books, err := model.GetBooks(dbx)
	if hutil.ReportError(writer, err) {
		return
	}

	hutil.RenderJSON(writer, books)
}

type reqNewBookAuthor struct {
	BookID   int `json:"bookID"`
	AuthorID int `json:"authorID"`
}

// NewBookAuthorHandler handles /book/author/new and registers a book author.
func NewBookAuthorHandler(
	ctx web.C, writer http.ResponseWriter, httpReq *http.Request,
	dbx *sqlx.DB, body []byte, _ model.User,
) {
	var req reqNewBookAuthor
	err := json.Unmarshal(body, &req)
	if hutil.ReportError(writer, err) {
		return
	}

	err = model.AddBookAuthor(dbx, req.BookID, req.AuthorID)
	if hutil.ReportError(writer, err) {
		return
	}

	hutil.ReportOK(writer)
}

type reqDeleteBookAuthor struct {
	BookID   int `json:"bookID"`
	AuthorID int `json:"authorID"`
}

// DeleteBookAuthorHandler handles /book/author/delete and unregisters the book author.
func DeleteBookAuthorHandler(
	ctx web.C, writer http.ResponseWriter, httpReq *http.Request,
	dbx *sqlx.DB, body []byte, _ model.User,
) {
	var req reqDeleteBookAuthor
	err := json.Unmarshal(body, &req)
	if hutil.ReportError(writer, err) {
		return
	}

	err = model.RemoveBookAuthor(dbx, req.BookID, req.AuthorID)
	if hutil.ReportError(writer, err) {
		return
	}

	hutil.ReportOK(writer)
}
