package main

// TODO: handle erors properly

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {

	// get the port heroku assignened for us
	port := os.Getenv("PORT")

	if port == "" { // ....altså hvis heroku ikke har gitt oss en port (DEBUG)
		port = "8080"
	}

	// set up default path
	http.HandleFunc("/", handleBadRequest)

	// set up request handler
	http.HandleFunc("/projectinfo/v1/github.com/", handleRequest)

	// start listening on port 8080
	fmt.Println("listening on port " + port + "...")
	err := http.ListenAndServe(":"+port, nil)

	// if error, panic
	if err != nil {
		panic(err)
	}
}

// handler for when invalid path is used
func handleBadRequest(res http.ResponseWriter, req *http.Request) {
	status := http.StatusNotFound
	http.Error(res, http.StatusText(status), status)
}

// handler for dealing with requests
func handleRequest(res http.ResponseWriter, req *http.Request) {

	// only GET is legal. only handle GET
	if req.Method == "GET" {

		// split URL to fetch information ([3:] are interesting)
		parts := strings.Split(req.URL.String(), "/")

		// hvis den MINST har user (4) og repo (5). hvis ikke, bad request
		if len(parts) > 5 {
			user := parts[4]
			repo := parts[5]

			// generate payload based on user, repo
			payload := generateResponsePayload(user, repo)

			// errorcheck payload and write
			if true {
				json.NewEncoder(res).Encode(payload)
			} else {
				// TODO denne kan være 404 eller 503?
				status := http.StatusServiceUnavailable
				http.Error(res, http.StatusText(status), status)
			}
		} else {
			handleBadRequest(res, req)
		}

	} else {
		// not GET. bad
		status := http.StatusNotImplemented
		http.Error(res, http.StatusText(status), status)
	}
}
