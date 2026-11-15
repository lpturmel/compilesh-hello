// A small web server for compile.sh end-to-end tests. It listens on $PORT and
// has endpoints that answer with chosen status codes, so the Requests page has
// 2xx, 4xx and 5xx traffic to show.
//
//	GET  /                  200, what the platform told the app
//	GET  /health            200
//	POST /items             201 with {"name": "..."}, 400 otherwise
//	GET  /bad-request       400
//	GET  /error             500
//	GET  /status/{code}     any status from 200 to 599
//	GET  /slow?ms=1500      200 after a delay (at most 10 s)
//	anything else           404
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var version = "v4"

const maxSlow = 10 * time.Second

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func problem(w http.ResponseWriter, status int, detail string) {
	writeJSON(w, status, map[string]any{
		"status": status,
		"error":  http.StatusText(status),
		"detail": detail,
	})
}

func newMux(port string) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello from compile.sh (%s)\nhost=%s port=%s service=%s deployment=%s uid=%d\n",
			version, r.Host, port, os.Getenv("COMPILESH_SERVICE_NAME"), os.Getenv("COMPILESH_DEPLOYMENT_ID"), os.Getuid())
	})

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "version": version})
	})

	mux.HandleFunc("POST /items", func(w http.ResponseWriter, r *http.Request) {
		var item struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<16)).Decode(&item); err != nil {
			problem(w, http.StatusBadRequest, "body must be JSON like {\"name\": \"...\"}")
			return
		}
		if item.Name == "" {
			problem(w, http.StatusBadRequest, "name is required")
			return
		}
		writeJSON(w, http.StatusCreated, map[string]string{"name": item.Name})
	})

	mux.HandleFunc("GET /bad-request", func(w http.ResponseWriter, r *http.Request) {
		problem(w, http.StatusBadRequest, "this endpoint always rejects the request")
	})

	mux.HandleFunc("GET /error", func(w http.ResponseWriter, r *http.Request) {
		problem(w, http.StatusInternalServerError, "simulated failure")
	})

	mux.HandleFunc("GET /status/{code}", func(w http.ResponseWriter, r *http.Request) {
		code, err := strconv.Atoi(r.PathValue("code"))
		if err != nil || code < 200 || code > 599 {
			problem(w, http.StatusBadRequest, "code must be a number from 200 to 599")
			return
		}
		writeJSON(w, code, map[string]any{"status": code, "text": http.StatusText(code)})
	})

	mux.HandleFunc("GET /slow", func(w http.ResponseWriter, r *http.Request) {
		ms, err := strconv.Atoi(r.URL.Query().Get("ms"))
		if err != nil || ms < 0 {
			ms = 1500
		}
		delay := min(time.Duration(ms)*time.Millisecond, maxSlow)
		select {
		case <-time.After(delay):
			writeJSON(w, http.StatusOK, map[string]any{"slept_ms": delay.Milliseconds()})
		case <-r.Context().Done():
		}
	})

	// Everything else, including wrong methods on known paths' siblings.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		problem(w, http.StatusNotFound, "no route for "+r.Method+" "+r.URL.Path)
	})

	return mux
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, newMux(port)))
}
