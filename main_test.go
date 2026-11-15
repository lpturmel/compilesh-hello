package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStatuses(t *testing.T) {
	server := httptest.NewServer(newMux("8080"))
	defer server.Close()

	cases := []struct {
		method, path, body string
		want               int
	}{
		{"GET", "/", "", 200},
		{"GET", "/health", "", 200},
		{"POST", "/items", `{"name":"widget"}`, 201},
		{"POST", "/items", `{}`, 400},
		{"POST", "/items", `not json`, 400},
		{"GET", "/bad-request", "", 400},
		{"GET", "/error", "", 500},
		{"GET", "/status/503", "", 503},
		{"GET", "/status/204", "", 204},
		{"GET", "/status/99", "", 400},
		{"GET", "/status/abc", "", 400},
		{"GET", "/slow?ms=10", "", 200},
		{"GET", "/nope", "", 404},
		{"GET", "/items/1", "", 404},
	}
	for _, c := range cases {
		request, _ := http.NewRequest(c.method, server.URL+c.path, strings.NewReader(c.body))
		response, err := http.DefaultClient.Do(request)
		if err != nil {
			t.Fatalf("%s %s: %v", c.method, c.path, err)
		}
		response.Body.Close()
		if response.StatusCode != c.want {
			t.Errorf("%s %s = %d, want %d", c.method, c.path, response.StatusCode, c.want)
		}
	}
}
