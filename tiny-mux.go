package tinymux

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
)

var methods = []string{"GET", "POST", "PATCH", "PUT", "DELETE"}

// keys type is unexported to prevent collisiosn with context keys
type keys int

const (
	paramsKey keys = iota
)

type radixTree struct {
	root *radixNode
}

type radixNode struct {
	// handler       http.Handler
	partial       string
	actualPattern string
	childs        map[string]*radixNode
	methods       map[string]http.Handler
}

type middleware func(http.Handler) http.Handler

func (m middleware) chainMiddleware(handler http.Handler) http.Handler {
	return m(handler)
}

func ChainMiddlewares(handler http.Handler, middlewares ...middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i].chainMiddleware(handler)
	}
	return handler
}

type TinyMux struct {
	radixTree   *radixTree
	middlewares []middleware
}

func (rt *radixTree) insert(method string, urlPattern string, handler http.Handler) {
	partialURL := partialURL(urlPattern)

	if rt.root == nil {
		rootNode := createRadixNode("/", nil)
		rt.root = rootNode
	}

	currNode := rt.root
	for i := 0; i < len(partialURL); i++ {
		partial := partialURL[i]
		childs := currNode.childs

		if strings.HasPrefix(partial, ":") {
			partial = "#"
		}

		child, childFound := childs[partial]

		if childFound {
			currNode = child
		} else {
			newChild := createRadixNode(partial, nil)
			currNode.childs[partial] = newChild
			currNode = newChild
		}
	}

	_, handlerFound := currNode.methods[method]

	if handlerFound {
		log.Fatal(fmt.Sprintf("two url has conflict with each other %s -- %s  %s", currNode.actualPattern, urlPattern, method))
	} else {
		currNode.methods[method] = handler
		currNode.actualPattern = urlPattern
	}
}

func createRadixNode(partial string, handler http.Handler) *radixNode {
	return &radixNode{
		// handler: handler,
		partial: partial,
		childs:  make(map[string]*radixNode),
		methods: make(map[string]http.Handler),
	}
}

func (rt *radixTree) search(urlPattern string) *radixNode {
	partialURL := partialURL(urlPattern)

	currNode := rt.root
	childFound := true
	for i := 0; i < len(partialURL); i++ {
		partial := partialURL[i]

		var child *radixNode
		childs := currNode.childs
		child, childFound = childs[partial]

		if childFound {
			currNode = child
		} else {
			child, childFound = childs["#"]
			if childFound {
				currNode = child
			}
		}

		if !childFound {
			return nil
		}
	}

	if childFound && len(currNode.methods) > 0 {
		return currNode
	}
	return nil
}

func methodNotAllowedHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("method not allowed"))
	})
}

func NewTinyMux() *TinyMux {
	radixTree := new(radixTree)
	tinyMux := TinyMux{
		radixTree: radixTree,
	}

	return &tinyMux
}
func (tm *TinyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// find corresponding handler, then call the matched handler ServeHTTP method
	url := r.URL.Path
	method := r.Method
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
	handlerNode := tm.radixTree.search(url)

	var handler http.Handler
	newR := r
	if handlerNode != nil {
		var ok bool
		handler, ok = handlerNode.methods[method]
		if !ok {
			handler = methodNotAllowedHandler()
		} else {
			newR = tm.readParamsValue(r, handlerNode)
		}
	} else {
		handler = http.NotFoundHandler()
	}

	handler = ChainMiddlewares(handler, tm.middlewares...)

	handler.ServeHTTP(w, newR)
}

func (tm *TinyMux) readParamsValue(r *http.Request, handlerNode *radixNode) *http.Request {

	actualURL := handlerNode.actualPattern
	url := r.URL.Path

	actualPartialURL := partialURL(actualURL)
	partialURL := partialURL(url)

	params := map[string]string{}
	for i, actualPartial := range actualPartialURL {

		if strings.HasPrefix(actualPartial, ":") {
			params[actualPartial[1:]] = partialURL[i]
		}
	}

	// https://stackoverflow.com/questions/40891345/fix-should-not-use-basic-type-string-as-key-in-context-withvalue-golint
	// TODO fix above issue
	rcopy := r.WithContext(context.WithValue(r.Context(), paramsKey, params))

	return rcopy

}

func Values(r http.Request) map[string]string {
	return r.Context().Value(paramsKey).(map[string]string)
}

func (tm *TinyMux) Handle(method string, urlPattern string, handler http.Handler) {
	if !existIn(method, methods) {
		log.Fatal("method is not valid", method)
	}
	// chain middlewares
	// handler = ChainMiddlewares(handler, tm.middlewares...)

	tm.radixTree.insert(method, urlPattern, handler)
}

func (tm *TinyMux) GET(urlPattern string, handler http.Handler) {
	tm.Handle("GET", urlPattern, handler)
}

func (tm *TinyMux) POST(urlPattern string, handler http.Handler) {
	tm.Handle("POST", urlPattern, handler)
}

func (tm *TinyMux) PATCH(urlPattern string, handler http.Handler) {
	tm.Handle("PATCH", urlPattern, handler)
}

func (tm *TinyMux) PUT(urlPattern string, handler http.Handler) {
	tm.Handle("PUT", urlPattern, handler)
}

func (tm *TinyMux) DELETE(urlPattern string, handler http.Handler) {
	tm.Handle("DELETE", urlPattern, handler)
}

func (tm *TinyMux) Use(middlewares ...middleware) {
	for _, handler := range middlewares {
		tm.middlewares = append(tm.middlewares, handler)
	}
}

func partialURL(urlPattern string) []string {

	if !strings.HasPrefix(urlPattern, "/") {
		panic("invalid urlPattern")
	}

	// urlPattern = strings.ReplaceAll(urlPattern, "/", "#/#")
	partialURL := strings.Split(urlPattern, "/")
	if partialURL[len(partialURL)-1] == "" {
		return partialURL[1 : len(partialURL)-1]
	}
	return partialURL[1:]
}

func existIn(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
