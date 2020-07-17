package hutil

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/akouryy/ss-api.pyon.app/src/model"
	"github.com/jmoiron/sqlx"
	"github.com/zenazn/goji/web"
	"golang.org/x/crypto/bcrypt"
)

func SetDBXMiddleware(dbx *sqlx.DB) func(*web.C, http.Handler) http.Handler {
	return func(ctx *web.C, httpHandler http.Handler) http.Handler {
		fn := func(writer http.ResponseWriter, httpReq *http.Request) {
			ctx.Env["dbx"] = dbx
			httpHandler.ServeHTTP(writer, httpReq)
		}

		return http.HandlerFunc(fn)
	}
}

func DBX(ctx *web.C) *sqlx.DB {
	return ctx.Env["dbx"].(*sqlx.DB)
}

func ReadWholeBody(httpReq *http.Request) ([]byte, error) {
	size, err := strconv.Atoi(httpReq.Header.Get("Content-Length"))
	if err != nil {
		return nil, err
	}
	body := make([]byte, size)
	size, err = httpReq.Body.Read(body)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return body, nil
}

func RenderJSON(w http.ResponseWriter, data interface{}) {
	j, err := json.Marshal(data)
	if ReportError(w, err) {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, string(j))
}

func ReportOK(w http.ResponseWriter) {
	fmt.Fprint(w, `{"ok":true}`)
}

func ReportError(w http.ResponseWriter, err error) bool {
	if err != nil {
		w.Header().Set("Content-Type", "application/json")

		escaped, _ := json.Marshal(err.Error())
		fmt.Fprintf(w, `{"error":%s}`, escaped)
		return true
	} else {
		return false
	}
}

type reqAuth struct {
	Nickname string `json:"authNickname"`
	Password string `json:"authPassword"`
}

func Authenticate(body []byte, dbx *sqlx.DB) (model.User, error) {
	var req reqAuth
	err := json.Unmarshal(body, &req)
	if err != nil {
		return model.User{}, err
	}

	var user model.User
	err = dbx.Get(&user, "SELECT * FROM users WHERE nickname = ?", req.Nickname)
	if err == sql.ErrNoRows {
		return model.User{}, errors.New("Wrong nickname or password")
	} else if err != nil {
		return model.User{}, err
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(req.Password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return model.User{}, errors.New("Wrong nickname or password")
	}
	if err != nil {
		log.Println(err)
		return model.User{}, err
	}

	return user, nil
}
