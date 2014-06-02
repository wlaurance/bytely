package main

import (
	"net/http"

	"../database"
)

func CreateUser(rw http.ResponseWriter, req *http.Request) {
	var data database.User

	err := ReadJSONRequest(rw, req, &data)
	if err != nil {
		WriteJSONResponse(rw, ResponseBadRequest)
		return
	}

	user := database.NewUser(data.Email, data.Password)
	err = user.Insert()
	if err != nil {
		WriteJSONResponse(rw, ResponseUserExists)
		return
	}

	// Build a map for the response.
	response := make(map[string]string)
	response["token"] = user.Token

	WriteJSONResponse(rw, DataResponse{http.StatusOK, response})
}
