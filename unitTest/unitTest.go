package unitTest

import (
	"fmt"
	"net/http"
)

func HelloWorld(r http.ResponseWriter, req *http.Request) {
	fmt.Fprint(r, "Hello World")
}

func main() {
	http.HandleFunc("/", HelloWorld)
	http.ListenAndServe(":3000", nil)
}