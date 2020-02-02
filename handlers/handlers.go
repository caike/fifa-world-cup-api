package handlers

import (
	"encoding/json"
	"fifa-heroku/data"
	"net/http"
	"time"
)

// RootHandler returns an empty body status code
func RootHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusNoContent)
}

// ListWinners returns winners from the list
func ListWinners(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	year := req.URL.Query().Get("year")
	if year == "" {
		winners, err := data.ListAllJSON()
		if err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		res.Write(winners)
	} else {
		filteredWinners, err := data.ListAllByYear(year)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)
			return
		}
		res.Write(filteredWinners)
	}
}

// AddNewWinner adds new winner to the list
func AddNewWinner(res http.ResponseWriter, req *http.Request) {
	accessToken := req.Header.Get("X-ACCESS-TOKEN")
	isTokenValid := data.IsAccessTokenValid(accessToken)
	if !isTokenValid {
		res.WriteHeader(http.StatusUnauthorized)
	} else {
		err := data.AddNewWinner(req.Body)
		if err != nil {
			res.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		res.WriteHeader(http.StatusCreated)
	}

}

// WinnersHandler is the dispatcher for all /winners URL
func WinnersHandler(res http.ResponseWriter, req *http.Request) {
	(res).Header().Set("Access-Control-Allow-Origin", "*")

	switch req.Method {
	case http.MethodGet:
		ListWinners(res, req)
	case http.MethodPost:
		AddNewWinner(res, req)
	default:
		res.WriteHeader(http.StatusMethodNotAllowed)
	}
}

type login struct {
	Email    string
	Password string
}

type user struct {
	Email string
	Token string
}

// LoginHandler validates user credentials
func LoginHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method == "OPTIONS" {
		res.Header().Add("Access-Control-Allow-Origin", "*")
		res.Header().Add("Access-Control-Allow-Headers", "Content-Type")
		res.WriteHeader(http.StatusOK)
		return
	}
	time.Sleep(3 * time.Second)
	res.Header().Add("Access-Control-Allow-Origin", "*")
	//res.Header().Add("Content-Type", "application/json")

	var logi login

	err := json.NewDecoder(req.Body).Decode(&logi)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	if logi.Email != "foo@bar.com" || logi.Password != "secret" {
		http.Error(res, "Invalid Credentials", http.StatusUnauthorized)
		return
	}
	user := user{Email: logi.Email, Token: "eyJ1eXAiOoJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJjcmVhdGVkIjoxNTgwNTY4MjgyMDAwLCJleHAiOjE1ODExNzMwODIwMDAsInBlc3NvYV9pZCI6NDQyMjEsInVzZXJuYW1lIjoiY2FybG9zQGlkb3B0ZXJsYWJzLmNvbS5iciJ9.6JyJvGRaT3X8fcF_HkDl7YFJPIr8ZqG8rBr5pbC1hBo"}
	js, err := json.Marshal(user)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Write(js)
}

/**
 curl -X POST "http://localhost:8080/login" \
 -H "accept: application/json" \
 -H "Content-Type: application/json" \
 -d "{\"email\":\"foo@bar.com\",\"password\":\"secret\"}"
**/
