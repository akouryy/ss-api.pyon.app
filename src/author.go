package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/zenazn/goji/web"
)

type Author struct {
	Id        int
	Name      string
	URL       string
	CreatedAt time.Time `db:"created_at"`
}

type reqNewAuthor struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func (author reqNewAuthor) Validate() error {
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

func AuthorsHandler(c web.C, writer http.ResponseWriter, httpReq *http.Request) {
	body, err := ReadWholeBody(httpReq)
	if ReportError(writer, err) {
		return
	}

	_, err = Authenticate(body)
	if ReportError(writer, err) {
		return
	}

	authors := []Author{}
	if ReportError(writer,
		dbx.Select(&authors, "SELECT * FROM authors ORDER BY created_at DESC LIMIT 10;"),
	) {
		return
	}

	RenderJSON(writer, authors)
}

func NewAuthorHandler(c web.C, writer http.ResponseWriter, httpReq *http.Request) {
	body, err := ReadWholeBody(httpReq)
	if ReportError(writer, err) {
		return
	}

	_, err = Authenticate(body)
	if ReportError(writer, err) {
		return
	}

	var req reqNewAuthor
	if ReportError(writer,
		json.Unmarshal(body, &req),
	) {
		return
	}

	if ReportError(writer,
		req.Validate(),
	) {
		return
	}

	_, err = dbx.NamedExec(
		"INSERT INTO authors(name, url, created_at) VALUES (:name, :url, NOW())",
		req,
	)
	if ReportError(writer, err) {
		return
	}

	ReportOK(writer)
}
