package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func ReportError(w http.ResponseWriter, err error) bool {
	if err != nil {
		escaped, _ := json.Marshal(err.Error())
		fmt.Fprintf(w, "{\"error\":%s}", escaped)
		return true
	} else {
		return false
	}
}

type reqAuth struct {
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

func Authenticate(httpReq *http.Request) (User, error) {
	var req reqAuth
	err := json.NewDecoder(httpReq.Body).Decode(&req)
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
