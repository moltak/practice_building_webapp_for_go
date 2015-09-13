package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/", HelloWorld)
	http.ListenAndServe(":3000", nil)
}

func HelloWorld(rw http.ResponseWriter, r *http.Request) {
	rw.Write([]byte("Hello World"))
}