package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_App(t *testing.T) {
	ts := httptest.NewServer(App())
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		t.Fatal(err)
	}

	exp := "Before...Hello World...After"

	if exp != string(body) {
		t.Fatalf("Expected %s got %s", exp, body)
	}
}

func Test_AppPost(t *testing.T) {
	ts := httptest.NewServer(App())
	defer ts.Close()

	bodyType := "string"
	res, _ := http.Post(ts.URL, bodyType, nil)
	body, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	exp := "Before...Hello Post...After"
	if exp != string(body) {
		t.Fatalf("Expected %s got %s", exp, body)
	}
}