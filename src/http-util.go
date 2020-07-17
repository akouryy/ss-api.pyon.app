package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

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

func Authenticate(body []byte) (User, error) {
	var req reqAuth
	err := json.Unmarshal(body, &req)
	if err != nil {
		return User{}, err
	}

	var user User
	err = dbx.Get(&user, "SELECT * FROM users WHERE nickname = ?", req.Nickname)
	if err == sql.ErrNoRows {
		return User{}, errors.New("Wrong nickname or password")
	} else if err != nil {
		return User{}, err
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(req.Password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return User{}, errors.New("Wrong nickname or password")
	}
	if err != nil {
		log.Println(err)
		return User{}, err
	}

	return user, nil
}
