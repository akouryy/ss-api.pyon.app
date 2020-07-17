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
	BookId int `json:"bookID"`
}

func BookHandler(
	ctx web.C, writer http.ResponseWriter, httpReq *http.Request,
	dbx *sqlx.DB, body []byte, _ model.User,
) {
	var req reqBook
	err := json.Unmarshal(body, &req)
	if hutil.ReportError(writer, err) {
		return
	}

	book, err := model.GetBook(dbx, req.BookId)
	if hutil.ReportError(writer, err) {
		return
	}

	hutil.RenderJSON(writer, book)
}

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
	BookId   int `json:"bookID"`
	AuthorId int `json:"authorID"`
}

func NewBookAuthorHandler(
	ctx web.C, writer http.ResponseWriter, httpReq *http.Request,
	dbx *sqlx.DB, body []byte, _ model.User,
) {
	var req reqNewBookAuthor
	err := json.Unmarshal(body, &req)
	if hutil.ReportError(writer, err) {
		return
	}

	err = model.AddBookAuthor(dbx, req.BookId, req.AuthorId)
	if hutil.ReportError(writer, err) {
		return
	}

	hutil.ReportOK(writer)
}

type reqDeleteBookAuthor struct {
	BookId   int `json:"bookID"`
	AuthorId int `json:"authorID"`
}

func DeleteBookAuthorHandler(
	ctx web.C, writer http.ResponseWriter, httpReq *http.Request,
	dbx *sqlx.DB, body []byte, _ model.User,
) {
	var req reqDeleteBookAuthor
	err := json.Unmarshal(body, &req)
	if hutil.ReportError(writer, err) {
		return
	}

	err = model.RemoveBookAuthor(dbx, req.BookId, req.AuthorId)
	if hutil.ReportError(writer, err) {
		return
	}

	hutil.ReportOK(writer)
}
