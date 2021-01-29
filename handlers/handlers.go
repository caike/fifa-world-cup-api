package handlers

import (
	"encoding/json"
	"fifa-heroku/data"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = os.Getenv("JWT_SECRET")

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

type login struct {
	Email    string
	Password string
}

type user struct {
	Token string
}
// // AddNewWinner adds new winner to the list
// func AddNewWinner(res http.ResponseWriter, req *http.Request) {
// 	accessToken := req.Header.Get("X-ACCESS-TOKEN")
// 	isTokenValid := data.IsAccessTokenValid(accessToken)
// 	if !isTokenValid {
// 		res.WriteHeader(http.StatusUnauthorized)
// 	} else {
// 		err := data.AddNewWinner(req.Body)
// 		if err != nil {
// 			res.WriteHeader(http.StatusUnprocessableEntity)
// 			return
// 		}
// 		res.WriteHeader(http.StatusCreated)
// 	}

// }

// WinnersHandler is the dispatcher for all /winners URL
func WinnersHandler(res http.ResponseWriter, req *http.Request) {
	(res).Header().Set("Access-Control-Allow-Origin", "*")

	var loggedUser user

	err := json.NewDecoder(req.Body).Decode(&loggedUser)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := jwt.Parse(loggedUser.Token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(jwtSecret), nil
	})


	if _, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		http.Error(res, "Invalid Credentials", http.StatusUnauthorized)
		return
	}

	switch req.Method {
	case http.MethodGet:
		ListWinners(res, req)
	// case http.MethodPost:
	// 	AddNewWinner(res, req)
	default:
		res.WriteHeader(http.StatusMethodNotAllowed)
	}
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
	res.Header().Add("Content-Type", "application/json")

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
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": logi.Email,
		"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil{
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	user := user{Token: tokenString}

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

 // read jwt token and then make authenticated
 // calls as so:

 curl -X GET "http://localhost:8080/winners" \
 -H "accept: application/json" \
 -H "Content-Type: application/json" \
 -d "{\"Token\":\"use-jwt-token-here\"}"
**/
