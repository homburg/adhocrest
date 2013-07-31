package main

import (
	"strings"
	"net/http"
	"log"
	"runtime"
	"fmt"
)

var address string

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	address = "localhost:31415"

	http.HandleFunc("/", index)

	log.Printf("Starting server (%s)...\n", address)
	http.ListenAndServe(address, nil)
}

var html string

func respond(w http.ResponseWriter, r *http.Request) {
	headers := []string{
		"Allow",
		"Access-Control-Allow-Methods",
	}

	for _, header := range headers {
		w.Header().Set(header, "POST,PUT,DELETE,GET,OPTIONS")
	}
	w.Header().Set("Access-Control-Allow-Origins", "*")
	fmt.Fprintln(w, r.Method)
}

func index(w http.ResponseWriter, r *http.Request) {
	switch strings.ToUpper(r.Method) {
		case "OPTIONS":
			// Cross-origin post
			log.Printf(r.Method)
			respond(w, r)
		case "POST":
			// create
			log.Printf(r.Method)
			respond(w, r)
		case "PUT":
			// update
			log.Printf(r.Method)
			respond(w, r)
		case "DELETE":
			log.Printf(r.Method)
			respond(w, r)
			// drop
		case "GET":
			// read
			log.Printf(r.Method)
			respond(w, r)
		case "HEAD":
			log.Println("Unsupported method")
			w.WriteHeader(http.StatusMethodNotAllowed)
		default:
			log.Printf("Unknown request method \"%s\"\n", r.Method)
			w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

