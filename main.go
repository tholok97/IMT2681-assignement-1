package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/projectinfo/v1/", hello)
	fmt.Println("listening...")
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		panic(err)
	}
}

func hello(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "hello, world")
}
