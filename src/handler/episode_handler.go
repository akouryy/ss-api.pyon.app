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

type reqEpisode struct {
	EpisodeID int `json:"episodeID"`
}

type reqNewEpisode struct {
	Title    string `json:"title"`
	BookID   int    `json:"bookID"`
	IndexNum int    `json:"index"`
}

func (episode *reqNewEpisode) validate() error {
	if episode.Title == "" {
		return errors.New("episode title must be nonempty")
	}
	if 100 < len(episode.Title) {
		return errors.New("episode title must be at most 100 characters")
	}
	if episode.IndexNum <= 0 {
		return errors.New("episode index must be at least 1")
	}
	return nil
}

// EpisodeHandler handles /episode/show and responds with an episode.
func EpisodeHandler(
	ctx web.C, writer http.ResponseWriter, httpReq *http.Request,
	dbx *sqlx.DB, body []byte, _ model.User,
) {
	var req reqEpisode
	err := json.Unmarshal(body, &req)
	if hutil.ReportError(writer, err) {
		return
	}

	bookEpi, err := model.GetEpisodeWithBook(dbx, req.EpisodeID)
	if hutil.ReportError(writer, err) {
		return
	}

	hutil.RenderJSON(writer, bookEpi)
}

// NewEpisodeHandler handles /episode/new and creates an episode.
func NewEpisodeHandler(
	ctx web.C, writer http.ResponseWriter, httpReq *http.Request,
	dbx *sqlx.DB, body []byte, _ model.User,
) {
	var req reqNewEpisode
	err := json.Unmarshal(body, &req)
	if hutil.ReportError(writer, err) {
		return
	}

	err = req.validate()
	if hutil.ReportError(writer, err) {
		return
	}

	err = model.CreateEpisode(dbx, req.Title, req.BookID, req.IndexNum)
	if hutil.ReportError(writer, err) {
		return
	}

	hutil.ReportOK(writer)
}
