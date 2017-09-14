package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {

	// set up default path
	http.HandleFunc("/", handleNotFound)

	// set up request handler
	http.HandleFunc("/projectinfo/v1/", handleRequest)

	// start listening on port 8080
	fmt.Println("listening...")
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)

	// if error, panic
	if err != nil {
		panic(err)
	}
}

// handler for when invalid path is used
func handleNotFound(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "404 - Not found")
}

// handler for dealing with requests
func handleRequest(res http.ResponseWriter, req *http.Request) {

	// only GET is legal. only handle GET
	if req.Method == "GET" {

		// split URL to fetch information ([3:] are interesting)
		parts := strings.Split(req.URL.String(), "/")

		// print them (FOR DEBUG PURPOSES)
		for _, p := range parts[3:] {
			fmt.Fprint(res, p, ", ")
		}
		fmt.Fprintln(res)

	} else {
		// not GET. bad
		fmt.Fprintln(res, "501 - Not implemented")
	}
}
