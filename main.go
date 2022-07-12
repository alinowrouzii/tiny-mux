package main

import (
	"fmt"
	"net/http"
	"strings"
)

func pingHandler(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("PONG"))
	return
}

func main() {

	// url := "hello/world/:ali/something/:go/"
	url := ""
	fmt.Println(normalizeUrl(url))

	// tmMux := modules.NewTinyMux()

	// tmMux.Handle("/hello", http.HandlerFunc(pingHandler)).Methods("GET")

	// http.ListenAndServe(":8000", tmMux)
}

func normalizeUrl(urlPattern string) string {

	partialUrl := strings.Split(urlPattern, "/")
	fmt.Println(partialUrl)

	for index, partial := range partialUrl {
		if strings.HasPrefix(partial, ":") {
			partialUrl[index] = "*"
		}
	}

	normalizedUrl := "/" + strings.Join(partialUrl, "/")

	if normalizedUrl == "//" {
		return "/"
	}

	return normalizedUrl
}
