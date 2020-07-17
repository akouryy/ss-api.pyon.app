package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/akouryy/ss-api.pyon.app/src/handler/hutil"
	"github.com/akouryy/ss-api.pyon.app/src/model"
	"github.com/jmoiron/sqlx"
	"github.com/zenazn/goji/web"
)

type reqNewAuthor struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func (author *reqNewAuthor) validate() error {
	if author.Name == "" {
		return errors.New("Author name must be nonempty.")
	}
	if 30 < len(author.Name) {
		return errors.New("Author name must be at most 30 characters.")
	}
	if author.URL == "" {
		return errors.New("Author URL must be nonempty.")
	}
	if 300 < len(author.URL) {
		return errors.New("Author name must be at most 300 characters.")
	}
	return nil
}

func AuthorsHandler(
	ctx web.C, writer http.ResponseWriter, httpReq *http.Request,
	dbx *sqlx.DB, body []byte, _ model.User,
) {
	authors, err := model.GetAuthors(dbx)
	if hutil.ReportError(writer, err) {
		return
	}

	hutil.RenderJSON(writer, authors)
}

func NewAuthorHandler(
	ctx web.C, writer http.ResponseWriter, httpReq *http.Request,
	dbx *sqlx.DB, body []byte, _ model.User,
) {
	var req reqNewAuthor
	err := json.Unmarshal(body, &req)
	if hutil.ReportError(writer, err) {
		return
	}

	err = req.validate()
	if hutil.ReportError(writer, err) {
		return
	}

	err = model.CreateAuthor(dbx, req.Name, req.URL)
	if hutil.ReportError(writer, err) {
		return
	}

	hutil.ReportOK(writer)
}
