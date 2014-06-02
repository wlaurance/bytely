package main

import (
	"net/http"

	"../database"
)

func CreateLink(rw http.ResponseWriter, req *http.Request) {
	var link database.Link

	token := GetAuthToken(req)
	if token == "" {
		WriteJSONResponse(rw, ResponseInvalidToken)
		return
	}

	err := ReadJSONRequest(rw, req, &link)
	if err != nil {
		WriteJSONResponse(rw, ResponseBadRequest)
		return
	}

	user, err := database.GetUserByToken(token)
	if err != nil || user == nil {
		WriteJSONResponse(rw, ResponseInvalidToken)
		return
	}

	link.UserId = user.Id
	err = link.Insert()
	if err != nil {
		WriteJSONResponse(rw, ResponseInternalServerError)
		return
	}

	response := make(map[string]string)
	response["hash"] = link.Hash

	WriteJSONResponse(rw, DataResponse{http.StatusCreated, response})
}

func GetLink(rw http.ResponseWriter, req *http.Request, hash string) {
	link, err := database.GetLink(hash)
	if err != nil {
		WriteJSONResponse(rw, ResponseInternalServerError)
		return
	}

	WriteJSONResponse(rw, DataResponse{http.StatusOK, link})
}

func GetLinks(rw http.ResponseWriter, req *http.Request) {
	token := GetAuthToken(req)
	if token == "" {
		WriteJSONResponse(rw, ResponseInvalidToken)
		return
	}

	user, err := database.GetUserByToken(token)
	if err != nil || user == nil {
		WriteJSONResponse(rw, ResponseInvalidToken)
		return
	}

	links, err := database.GetLinksForUser(user.Id)
	if err != nil {
		WriteJSONResponse(rw, ResponseInternalServerError)
		return
	}

	WriteJSONResponse(rw, DataResponse{http.StatusOK, links})
}
