package main

import (
	"fmt"
	"net/http"
	"tiny-mux/modules"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println(modules.Values(*r))

	w.Write([]byte("PONG"))
	return
}

func main() {

	tmMux := modules.NewTinyMux()

	// tmMux.Handle("/hello/bar", http.HandlerFunc(pingHandler))
	// tmMux.Handle("/:hello/bar", http.HandlerFunc(pingHandler))
	// tmMux.Handle("/hello/world/", http.HandlerFunc(pingHandler))
	tmMux.Handle("GET", "/hello/world/:baz/bar", http.HandlerFunc(pingHandler))
	tmMux.Handle("GET", "/hello/world/bazz/bar", http.HandlerFunc(pingHandler))

	tmMux.Handle("POST", "/hello/world/:baz", http.HandlerFunc(pingHandler))
	// // conflict handling... goroutine must panic in below state
	// tmMux.Handle("/hello/world/:baz/", http.HandlerFunc(pingHandler))

	// handle this conflict:
	// tmMux.Handle("/:bar/", http.HandlerFunc(pingHandler))
	// tmMux.Handle("/:foo/bar/baz", http.HandlerFunc(pingHandler))
	// tmMux.Handle("/:foo/bar/bar", http.HandlerFunc(pingHandler))
	// tmMux.Handle("/:baz/bar/baz", http.HandlerFunc(pingHandler))

	tmMux.Handle("GET", "/:hello/bar/:id", http.HandlerFunc(pingHandler))
	// tmMux.Handle("/:hello/bar/baz", http.HandlerFunc(pingHandler))

	http.ListenAndServe(":8000", tmMux)
}
