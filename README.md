# FIFA API

Simple web API returning a list of FIFA World Cup winners.
Login handler expects `foo@bar.com / secret`

## Running the app

The app expects ENVs to be set for `JWT_SECRET` and `PORT`.
You can run it like so:

`JWT_SECRET=mysecret PORT=8080 go run main.go`

## Usage

```
 curl -X POST "http://localhost:8080/login" \
 -H "accept: application/json" \
 -H "Content-Type: application/json" \
 -d "{\"email\":\"foo@bar.com\",\"password\":\"secret\"}"
```

Read jwt token and then pass the jwt token in JSON like so:

```
 curl -X GET "http://localhost:8080/winners" \
 -H "accept: application/json" \
 -H "Content-Type: application/json" \
 -d "{\"Token\":\"use-jwt-token-here\"}"
```
