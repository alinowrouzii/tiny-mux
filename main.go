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
	tmMux.Handle("/hello/", http.HandlerFunc(pingHandler))
	// .Methods("GET")

	http.ListenAndServe(":8000", tmMux)
}

// func normalizeUrl(urlPattern string) string {

// 	partialUrl := strings.Split(urlPattern, "/")
// 	// fmt.Println(partialUrl)

// 	for index, partial := range partialUrl {
// 		if strings.HasPrefix(partial, ":") {
// 			partialUrl[index] = "*"
// 		}
// 	}

// 	normalizedUrl := "/" + strings.Join(partialUrl, "/")

// 	if normalizedUrl == "//" {
// 		return "/"
// 	}

// 	return normalizedUrl
// }

func partialUrl(urlPattern string) []string {
	if !strings.HasPrefix(urlPattern, "/") {
		panic("invalid urlPattern")
	}

	urlPattern = strings.ReplaceAll(urlPattern, "/", "#/#")

	partialUrl := strings.Split(urlPattern, "#")

	return partialUrl
}
