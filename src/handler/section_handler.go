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

type reqNewSection struct {
	Content   string `json:"content"`
	EpisodeId int    `json:"episodeID"`
	IndexNum  int    `json:"index"`
}

func (section *reqNewSection) validate() error {
	if section.Content == "" {
		return errors.New("Section content must be nonempty.")
	}
	if section.IndexNum <= 0 {
		return errors.New("Section index must be at least 1.")
	}
	return nil
}

func NewSectionHandler(
	ctx web.C, writer http.ResponseWriter, httpReq *http.Request,
	dbx *sqlx.DB, body []byte, _ model.User,
) {
	var req reqNewSection
	err := json.Unmarshal(body, &req)
	if hutil.ReportError(writer, err) {
		return
	}

	err = req.validate()
	if hutil.ReportError(writer, err) {
		return
	}

	err = model.CreateSection(dbx, req.Content, req.EpisodeId, req.IndexNum)
	if hutil.ReportError(writer, err) {
		return
	}

	hutil.ReportOK(writer)
}
