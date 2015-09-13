package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_HelloWorld(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com/foo", nil)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	HelloWorld(res, req)

	exp := "Hello World"
	act := res.Body.String()
	if exp != act {
		t.Fatalf("Expected %s gog %s", exp, act)
	}
}

func Test_HelloPost(t *testing.T) {
	req, err := http.NewRequest("POST", "http://example.com/foo", nil)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	HelloWorld(res, req)

	exp := "Hello Post"
	act := res.Body.String()
	if exp != act {
		t.Fatalf("Expected %s gog %s", exp, act)
	}
}