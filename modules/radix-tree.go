package modules

import (
	"net/http"
	"strings"
)

var Shit = 50

type radixTree struct {
	root *radixNode
}

type radixNode struct {
	handler http.Handler
	partial string
	childs  map[string]*radixNode
}

func (rt *radixTree) insert(pattern string, handler http.Handler) {

}

func (rt *radixTree) search(url string) http.Handler {

	return nil
}

type TinyMux struct {
	handlers map[string]*MuxHandler
}

type MuxHandler struct {
	handler           http.Handler
	methods           []string
	originPattern     string
	normalizedPattern string
}

func NewTinyMux() *TinyMux {
	handlers := make(map[string]*MuxHandler)
	tinyMux := TinyMux{
		handlers: handlers,
	}

	return &tinyMux
}
func (tm *TinyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// find corresponding handler, then call the matched handler ServeHTTP method
	url := r.URL.Path
	normalizedUrl := normalizeUrl(url)

	// This is where we need radixTree.
	// Forinstance user request for /foo/bar
	// AND a handler is registered for /foo/:bar
	// the radix tree should be able to iterate over tree.
	// Like: looking for foo --> found
	// search foo has a child or not --> foo has a child.
	// coressponding to previous assumption that we have
	// aready registerd /foo/:bar so the child of foo
	// is :bar (wild card). Therfore the pattern matched successfully
	// And  corresponding handler can be called.
	muxHandler, ok := tm.handlers[normalizedUrl]

	if !ok {
		http.NotFoundHandler().ServeHTTP(w, r)
		return
	}

	// check method is allowed
	method := r.Method
	if !existIn(method, muxHandler.methods) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(405)
		w.Write([]byte("method not allowed"))
		return
	}

	muxHandler.handler.ServeHTTP(w, r)
}

func (tm *TinyMux) Handle(urlPattern string, handler http.Handler) *MuxHandler {

	normalizedUrl := normalizeUrl(urlPattern)

	muxHandler := MuxHandler{
		handler:           handler,
		originPattern:     urlPattern,
		normalizedPattern: normalizedUrl,
	}

	tm.handlers[normalizedUrl] = &muxHandler

	return &muxHandler
}

func (muxHandler *MuxHandler) Methods(methods ...string) {
	muxHandler.methods = append(muxHandler.methods, methods...)
}

func existIn(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func normalizeUrl(urlPattern string) string {

	partialUrl := strings.Split(urlPattern, "/")
	// fmt.Println(partialUrl)

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

func partialUrl(urlPattern string) []string {

	if !strings.HasPrefix(urlPattern, "/") {
		panic("invalid urlPattern")
	}

	urlPattern = strings.ReplaceAll(urlPattern, "/", "#/#")
	partialUrl := strings.Split(urlPattern, "#")

	return partialUrl
}
