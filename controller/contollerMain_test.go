package main

import (
	"net/http"
	"testing"
	"net/http/httptest"
	"io/ioutil"

	"gopkg.in/unrolled/render.v1"
)

func Test_ControllerMainPost(t *testing.T) {
	c := &MyController{Render: render.New(render.Options{})}
	ts := httptest.NewServer(c.Action(c.Index))
	defer ts.Close()

	res, _ := http.Post(ts.URL, "string", nil)
	body, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	exp := "{\"Hello\":\"POST\"}"
	if exp != string(body) {
		t.Fatalf("Expected %s got %s", exp, body)
	}
}

func Test_ControllerMainGet(t *testing.T) {
	c := &MyController{Render: render.New(render.Options{})}
	ts := httptest.NewServer(c.Action(c.Index))
	defer ts.Close()

	res, _ := http.Get(ts.URL)
	body, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	exp := "{\"Hello\":\"GET\"}"
	if exp != string(body) {
		t.Fatalf("Expected %s got %s", exp, body)
	}
}