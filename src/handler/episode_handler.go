package handler

import (
	"encoding/json"
	"net/http"

	"github.com/akouryy/ss-api.pyon.app/src/handler/hutil"
	"github.com/akouryy/ss-api.pyon.app/src/model"
	"github.com/jmoiron/sqlx"
	"github.com/zenazn/goji/web"
)

type reqEpisode struct {
	EpisodeId int `json:"episodeID"`
}

func EpisodeHandler(
	ctx web.C, writer http.ResponseWriter, httpReq *http.Request,
	dbx *sqlx.DB, body []byte, _ model.User,
) {
	var req reqEpisode
	err := json.Unmarshal(body, &req)
	if hutil.ReportError(writer, err) {
		return
	}

	book, err := model.GetEpisode(dbx, req.EpisodeId)
	if hutil.ReportError(writer, err) {
		return
	}

	hutil.RenderJSON(writer, book)
}
