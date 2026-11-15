// A minimal web server for compile.sh end-to-end tests: it listens on $PORT
// and answers with what the platform told it.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var version = "v2"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello from compile.sh (%s)\nhost=%s port=%s service=%s deployment=%s uid=%d\n",
			version, r.Host, port, os.Getenv("COMPILESH_SERVICE_NAME"), os.Getenv("COMPILESH_DEPLOYMENT_ID"), os.Getuid())
	})
	log.Printf("listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
