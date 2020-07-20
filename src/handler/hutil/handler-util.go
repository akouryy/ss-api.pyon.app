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

var allowedOrigins = map[string]struct{}{
	"http://localhost:3000": struct{}{},
	"http://localhost:3001": struct{}{},
	"https://ss.pyon.app":   struct{}{},
}

// CORSMiddleware is a Goji middleware
// adding headers about Cross-Origin Resource Sharing.
func CORSMiddleware() func(*web.C, http.Handler) http.Handler {
	return func(ctx *web.C, httpHandler http.Handler) http.Handler {
		fn := func(writer http.ResponseWriter, httpReq *http.Request) {
			origin := httpReq.Header.Get("Origin")
			if _, ok := allowedOrigins[origin]; ok {
				writer.Header().Set("Access-Control-Allow-Origin", origin)
				writer.Header().Set("Access-Control-Allow-Headers", "*")
				writer.Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
			}

			httpHandler.ServeHTTP(writer, httpReq)
		}

		return http.HandlerFunc(fn)
	}
}

// ReadWholeBody stores the entire request body into a byte array.
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

// RenderJSON stringifies `data` as JSON and writes it to the response.
func RenderJSON(w http.ResponseWriter, data interface{}) {
	j, err := json.Marshal(data)
	if ReportError(w, err) {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, string(j))
}

// ReportOK writes an nullary successful JSON string to the response.
func ReportOK(w http.ResponseWriter) {
	fmt.Fprint(w, `{"ok":true}`)
}

// ReportError writes an JSON string containing the error message to the response.
func ReportError(w http.ResponseWriter, err error) bool {
	if err != nil {
		w.Header().Set("Content-Type", "application/json")

		escaped, _ := json.Marshal(err.Error())
		fmt.Fprintf(w, `{"error":%s}`, escaped)
		return true
	}
	return false
}

type reqAuth struct {
	Nickname string `json:"authNickname"`
	Password string `json:"authPassword"`
}

func (req *reqAuth) validate() error {
	if req.Nickname == "" {
		return errors.New("authenticate nickname must be nonempty")
	}
	if req.Password == "" {
		return errors.New("authenticate password must be nonempty")
	}
	return nil
}

// Authenticate confirms that the request is created by a registered user
// and identifies her/him.
func Authenticate(body []byte, dbx *sqlx.DB) (model.User, error) {
	var req reqAuth
	err := json.Unmarshal(body, &req)
	if err != nil {
		return model.User{}, err
	}

	err = req.validate()
	if err != nil {
		return model.User{}, err
	}

	var user model.User
	err = dbx.Get(&user, "SELECT * FROM users WHERE nickname = ?", req.Nickname)
	if err == sql.ErrNoRows {
		return model.User{}, errors.New("wrong nickname or password")
	} else if err != nil {
		return model.User{}, err
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(req.Password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return model.User{}, errors.New("wrong nickname or password")
	}
	if err != nil {
		log.Println(err)
		return model.User{}, err
	}

	return user, nil
}

// Wrap wraps a rich handler into a raw Goji handler.
func Wrap(
	dbx *sqlx.DB, raw func(web.C, http.ResponseWriter, *http.Request, *sqlx.DB, []byte, model.User),
) func(web.C, http.ResponseWriter, *http.Request) {
	return func(ctx web.C, writer http.ResponseWriter, httpReq *http.Request) {
		body, err := ReadWholeBody(httpReq)
		if ReportError(writer, err) {
			return
		}

		user, err := Authenticate(body, dbx)
		if ReportError(writer, err) {
			return
		}

		raw(ctx, writer, httpReq, dbx, body, user)
	}
}
