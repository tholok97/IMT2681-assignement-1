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

	if port == "" { // ....if heroku didn't give us a port (DEBUG)
		port = "8080"
	}

	// set up default path
	http.HandleFunc("/", handleBadRequest)

	// set up request handler
	http.HandleFunc("/projectinfo/v1/github.com/", handleRequest)

	// start listening on port 8080
	fmt.Println("Listening on port " + port + "...")
	err := http.ListenAndServe(":"+port, nil)

	// if error, panic
	if err != nil {
		panic(err)
	}
}

// handler for when invalid path is used
func handleBadRequest(res http.ResponseWriter, req *http.Request) {
	status := http.StatusBadRequest
	http.Error(res, http.StatusText(status), status)
}

// handler for dealing with requests
func handleRequest(res http.ResponseWriter, req *http.Request) {

	// only GET is legal. only handle GET
	if req.Method == "GET" {

		// split URL to fetch information ([3:] are interesting)
		parts := strings.Split(req.URL.String(), "/")

		// if it AT LEAST has user(4) and repo(5), do stuff. if not: bad request
		if len(parts) > 5 {
			user := parts[4]
			repo := parts[5]

			// generate payload based on user, repo
			payload, pErr := generateResponsePayload(user, repo)

			// errorcheck payload and write
			if pErr == nil {
				http.Header.Add(res.Header(),
					"content-type", "application/json")
				json.NewEncoder(res).Encode(payload)
			} else {

				// NOTE:	Need to show different status depending on error.
				//			The following is a BAD solution to the problem, very
				// TODO		"ad-hoc". Needs improvenemt
				var status int
				if pErr.Error() == "unexpected end of JSON input" {
					status = http.StatusBadRequest
				} else {
					status = http.StatusServiceUnavailable
				}
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
