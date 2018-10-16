package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		name := request.FormValue("name")
		fmt.Fprintf(writer, "<h1>hello world %s <h1>", name)
	})
	http.ListenAndServe(":8080", nil)
}
