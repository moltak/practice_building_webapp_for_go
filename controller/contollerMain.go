package main

import (
	"net/http"

	"gopkg.in/unrolled/render.v1"
)


type Action func(res http.ResponseWriter, req *http.Request) error

type AppController struct {}

func (c *AppController) Action(a Action) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		if err := a(res, req); err != nil {
			http.Error(res, err.Error(), 500)
		}
	})
}

type MyController struct {
	AppController
	*render.Render
}

func (c *MyController) Index(res http.ResponseWriter, req *http.Request) error {
	c.JSON(res, 200, map[string]string{"Hello":"Json"})
	return nil
}

func main() {
	c := &MyController{Render: render.New(render.Options{})}
	http.ListenAndServe(":3000", c.Action(c.Index))
}
