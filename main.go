package main

// TODO: handle erors properly

import (
	// "encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type Message struct {
	Name string
}

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
	res.WriteHeader(http.StatusBadRequest)
	fmt.Fprintln(res, "400 - Bad request")
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

			// generate payload based on user, repo and write to response
			fmt.Fprintln(res, generateResponsePayload(user, repo))
		} else {
			handleBadRequest(res, req)
		}

	} else {
		// not GET. bad
		res.WriteHeader(http.StatusNotImplemented)
		fmt.Fprintln(res, "501 - Not implemented")
	}
}

// send response to client
func generateResponsePayload(user, repo string) string {
	return "{\n\t\"project\": \"github.com/" + user + "/" + repo + "\",\n\t\"" + user + "\": \"apache\"\n}"
}
