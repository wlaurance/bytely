package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Responder assumes that any struct which can return a status code is a
// valid API response.
type Responder interface {
	StatusCode() int
}

// StatusResponse is for basic responses which consist only of an HTTP
// status code and a message.
type StatusResponse struct {
	status  int
	Message string `json:"message"`
}

func (this StatusResponse) StatusCode() int {
	return this.status
}

// 200 HTTP status codes.
var ResponseOK StatusResponse = StatusResponse{http.StatusOK, ""}
var ResponseCreated StatusResponse = StatusResponse{http.StatusCreated, ""}

// 400 HTTP status codes.
var ResponseUserExists StatusResponse = StatusResponse{
	http.StatusConflict,
	"A user already exists with that email.",
}
var ResponseUserDoesNotExist StatusResponse = StatusResponse{
	http.StatusNotFound,
	"A user with that email could not be located.",
}
var ResponseInvalidToken StatusResponse = StatusResponse{
	http.StatusUnauthorized,
	"Provided token is invalid.",
}
var ResponseBadRequest StatusResponse = StatusResponse{
	http.StatusBadRequest,
	"Request was malformed.",
}

// 500 HTTP status codes.
var ResponseInternalServerError StatusResponse = StatusResponse{
	http.StatusInternalServerError,
	"",
}

// DataResponse is for responses which are actually expected to return
// serialized database objects.
type DataResponse struct {
	status int
	data   interface{}
}

func (this DataResponse) StatusCode() int {
	return this.status
}

// GetAuthToken just wraps the extraction of an API token from a request's headers.
func GetAuthToken(req *http.Request) string {
	return req.Header.Get("X-Access-Token")
}

// SetCORSOptions adds appropriate access headers to each response for cross origin
// resource sharing.
func SetCORSOptions(rw http.ResponseWriter) {
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Access-Control-Allow-Headers", "X-Access-Token, Content-Type")
	rw.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,PATCH,DELETE,OPTIONS")
}

func ReadJSONRequest(rw http.ResponseWriter, req *http.Request, body interface{}) error {
	rawBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		http.Error(rw, "An unknown error occurred.", http.StatusInternalServerError)
		return err
	}
	err = json.Unmarshal(rawBody, body)
	if err != nil {
		WriteJSONResponse(rw, ResponseBadRequest)
		return err
	}

	return nil
}

func WriteJSONResponse(rw http.ResponseWriter, res Responder) error {
	var data interface{}

	switch response := res.(type) {
	case StatusResponse:
		data = response
	case DataResponse:
		data = response.data
	}

	response, err := json.Marshal(data)
	if err != nil {
		http.Error(rw, "", http.StatusInternalServerError)
		return err
	}

	rw.Header().Set("Content-Type", "application/json")
	SetCORSOptions(rw)
	rw.WriteHeader(res.StatusCode())
	rw.Write(response)

	return nil
}
