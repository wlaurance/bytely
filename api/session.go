package main

import (
	"database/sql"
	"net/http"
	"net/url"

	"../database"
)

// APILogin extracts the email and password fields from a GET request
// URL and attempts to retrieve the user with those credentials. If
// successful, the user's API token is returned and saved in the session store.
func GetToken(rw http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {
		WriteJSONResponse(rw, ResponseBadRequest)
		return
	}

	email := req.Form.Get("email")
	password := req.Form.Get("password")
	if email == "" || password == "" {
		WriteJSONResponse(rw, ResponseBadRequest)
		return
	}

	// Escape the @ sign in the email address.
	email, err = url.QueryUnescape(email)
	if err != nil {
		WriteJSONResponse(rw, ResponseBadRequest)
		return
	}

	user, err := database.GetUser(email, password)
	if err != nil {
		if err == sql.ErrNoRows {
			WriteJSONResponse(rw, ResponseUserDoesNotExist)
		} else {
			WriteJSONResponse(rw, ResponseInternalServerError)
		}
		return
	}

	response := make(map[string]string)
	response["token"] = user.Token

	WriteJSONResponse(rw, DataResponse{http.StatusOK, response})
}
