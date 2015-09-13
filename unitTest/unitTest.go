package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", HelloWorld)
	http.ListenAndServe(":3000", nil)
}

func HelloWorld(res http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		fmt.Fprint(res, "Hello World")
	} else {
		fmt.Fprint(res, "Hello Post")
	}
}