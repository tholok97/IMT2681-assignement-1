package main

// TODO: handle erors properly

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

			// generate payload based on user, repo
			payload := generateResponsePayload(user, repo)

			// errorcheck payload and write
			if payload != nil {

				fmt.Fprintln(res, payload)
			} else {
				// TODO denne kan være 404 eller 503?
				res.WriteHeader(http.StatusServiceUnavailable)
				fmt.Fprintln(res, "503 - Service unavailable")
			}
		} else {
			handleBadRequest(res, req)
		}

	} else {
		// not GET. bad
		res.WriteHeader(http.StatusNotImplemented)
		fmt.Fprintln(res, "501 - Not implemented")
	}
}

// generate payload by requesting github for the info we need and then basing
// the payload off the resopnse
func generateResponsePayload(user, repo string) []byte {

	// make request
	resp, err := http.Get("https://api.github.com/repos/" + user + "/" + repo)
	defer resp.Body.Close() // we need to close it when we're done

	// error cheks:

	if err != nil {
		return nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	// get body as bytes, and return
	b, err := ioutil.ReadAll(resp.Body)

	dec := json.NewDecoder(resp.Body)

	var s string
	dec.Decode(s)
	fmt.Println(s)

	return b

	//return "{\n\t\"project\": \"github.com/" + user + "/" + repo + "\",\n\t\"" + user + "\": \"apache\"\n}"
}
