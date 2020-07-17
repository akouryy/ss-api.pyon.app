package handler

import (
	"net/http"

	"github.com/akouryy/ss-api.pyon.app/src/handler/hutil"
	"github.com/akouryy/ss-api.pyon.app/src/model"
	"github.com/zenazn/goji/web"
)

func BooksHandler(ctx web.C, writer http.ResponseWriter, httpReq *http.Request) {
	dbx := hutil.DBX(&ctx)

	body, err := hutil.ReadWholeBody(httpReq)
	if hutil.ReportError(writer, err) {
		return
	}

	_, err = hutil.Authenticate(body, dbx)
	if hutil.ReportError(writer, err) {
		return
	}

	books, err := model.GetBooks(dbx)
	if hutil.ReportError(writer, err) {
		return
	}

	hutil.RenderJSON(writer, books)
}
