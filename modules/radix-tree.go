package modules

import (
	"fmt"
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

func (rt *radixTree) insert(urlPattern string, handler http.Handler) {
	partialUrl := partialUrl(urlPattern)
	fmt.Println(partialUrl, len(partialUrl), "h")

	if rt.root == nil {
		fmt.Println("null")
		rootNode := createRadixNode(partialUrl[0], nil)
		rt.root = rootNode
	}

	currNode := rt.root
	for i := 0; i < len(partialUrl)-1; i++ {
		partial := partialUrl[i+1]
		childs := currNode.childs
		fmt.Println("partial-->", partial)

		child, childFound := childs[partial]

		if childFound {
			fmt.Println("found")
			currNode = child
		} else {
			fmt.Println("not found")
			newChild := createRadixNode(partial, nil)
			currNode.childs[partial] = newChild
		}
	}

	if currNode.handler == nil {
		currNode.handler = handler
	} else {
		panic("url has conflict with another registered handler")
	}

	// fmt.Println("hahah")
}

func createRadixNode(partial string, handler http.Handler) *radixNode {
	return &radixNode{
		handler: handler,
		partial: partial,
		childs:  make(map[string]*radixNode),
	}
}

func (rt *radixTree) search(url string) http.Handler {
	return nil
}

type TinyMux struct {
	radixTree *radixTree
}

// type MuxHandler struct {
// 	handler           http.Handler
// 	methods           []string
// 	originPattern     string
// 	normalizedPattern string
// }

func NewTinyMux() *TinyMux {
	radixTree := new(radixTree)
	tinyMux := TinyMux{
		radixTree,
	}

	return &tinyMux
}
func (tm *TinyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// find corresponding handler, then call the matched handler ServeHTTP method
	url := r.URL.Path
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
	handler := tm.radixTree.search(url)

	handler.ServeHTTP(w, r)
}

func (tm *TinyMux) Handle(urlPattern string, handler http.Handler) {
	tm.radixTree.insert(urlPattern, handler)
}

// func (muxHandler *MuxHandler) Methods(methods ...string) {
// 	muxHandler.methods = append(muxHandler.methods, methods...)
// }

func existIn(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func partialUrl(urlPattern string) []string {

	if !strings.HasPrefix(urlPattern, "/") {
		panic("invalid urlPattern")
	}

	urlPattern = strings.ReplaceAll(urlPattern, "/", "#/#")
	partialUrl := strings.Split(urlPattern, "#")

	return partialUrl[1 : len(partialUrl)-1]
}
