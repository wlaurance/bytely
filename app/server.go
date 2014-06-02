package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"../common"
	"../database"
)

var API_ROOT = "http://localhost:1337"

func main() {
	// Create a multiplexer to serve up the app from this directory.
	fileMux := http.NewServeMux()
	fileMux.Handle("/", http.FileServer(http.Dir(".")))

	// Any requests that match the format API_ROOT/[A-z0-9]{6} will be treated
	// as shortened links and handled by attempting to retrieve the link's
	// original URL and redirect the user there. All other requests will be
	// forwarded to the file system.
	router := common.NewRouter()
	router.Route("/([A-z0-9]{6})", RedirectHandler)
	router.Route("(.*)", fileMux.ServeHTTP)

	err := http.ListenAndServe(":7331", router)
	if err != nil {
		log.Fatal("This web service sucks")
	}
}

// RedirectHandler receives the link hash that was extracted by the router and
// uses this to grab the full link object from the API. If successful, the
// user can be safely redirected to the expanded URL.
func RedirectHandler(rw http.ResponseWriter, req *http.Request, hash string) {
	link, err := database.GetLink(hash)
	if err != nil {
		log.Println(err)
		return
	}

	// Redirect the user to the page and finish processing.
	http.Redirect(rw, req, link.OriginalURL, http.StatusMovedPermanently)

	// Increase the mobile hit count if a mobile browser is detected.
	userAgent := req.Header.Get("User-Agent")
	if strings.Contains(userAgent, "Mobile") {
		link.MobileHits++
	}

	link.Hits++
	link.LastHit = time.Now().Format(time.RFC822)

	err = link.Save()
	if err != nil {
		log.Println("wieners")
		return
	}
}
