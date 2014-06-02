package main

import (
	"fmt"
	"log"
	"net/http"

	"../httputils"
)

// Maps each API route to a handler. The routes with regexes capture a chunk
// of the url and pass it to its respective handler as a string. Handlers
// are found below.
func main() {
	router := httputils.NewRouter()

	router.Route("/", RootHandler)
	router.Route("/users", UserHandler)
	router.Route("/links", LinkHandler)
	router.Route("/links/(.*)", LinkIdHandler)
	router.Route("/token", TokenHandler)
	router.Route("(.*)", NotFoundHandler)

	err := http.ListenAndServe(":1337", router)
	if err != nil {
		log.Fatal("Couldn't spin up da server ;_;", err)
	}
}

func RootHandler(rw http.ResponseWriter, req *http.Request) {
	WriteJSONResponse(rw, ResponseOK)
}

func UserHandler(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "POST":
		CreateUser(rw, req)
	case "OPTIONS":
		WriteJSONResponse(rw, ResponseOK)
	default:
		NotFoundHandler(rw, req)
	}
}

func LinkHandler(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		GetLinks(rw, req)
	case "POST":
		CreateLink(rw, req)
	case "OPTIONS":
		WriteJSONResponse(rw, ResponseOK)
	default:
		NotFoundHandler(rw, req)
	}
}

func LinkIdHandler(rw http.ResponseWriter, req *http.Request, param string) {
	switch req.Method {
	case "GET":
		GetLink(rw, req, param)
	case "PUT":
		fmt.Println("/links/:id PATCH")
	case "OPTIONS":
		WriteJSONResponse(rw, ResponseOK)
	default:
		NotFoundHandler(rw, req)
	}
}

func TokenHandler(rw http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		GetToken(rw, req)
	case "OPTIONS":
		WriteJSONResponse(rw, ResponseOK)
	default:
		NotFoundHandler(rw, req)
	}
}

func NotFoundHandler(rw http.ResponseWriter, req *http.Request) {
	http.Error(rw, "The requested URL was not found.", http.StatusNotFound)
}
