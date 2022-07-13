package main

import (
	"net/http"
	"strings"
	"tiny-mux/modules"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("PONG"))
	return
}

func main() {

	tmMux := modules.NewTinyMux()

	tmMux.Handle("/hello/", http.HandlerFunc(pingHandler))
	tmMux.Handle("/hello/world/", http.HandlerFunc(pingHandler))
	tmMux.Handle("/hello/world/:baz/bar", http.HandlerFunc(pingHandler))

	// conflict handling... goroutine must panic in below state
	tmMux.Handle("/hello/world/:baz", http.HandlerFunc(pingHandler))
	tmMux.Handle("/hello/world/:baz/", http.HandlerFunc(pingHandler))

	// handle this conflict:
	// tmMux.Handle("/:bar/", http.HandlerFunc(pingHandler))
	// tmMux.Handle("/:foo", http.HandlerFunc(pingHandler))
	// tmMux.Handle("/:baz", http.HandlerFunc(pingHandler))

	

	http.ListenAndServe(":8000", tmMux)
}


func partialUrl(urlPattern string) []string {
	if !strings.HasPrefix(urlPattern, "/") {
		panic("invalid urlPattern")
	}

	urlPattern = strings.ReplaceAll(urlPattern, "/", "#/#")

	partialUrl := strings.Split(urlPattern, "#")

	return partialUrl
}
