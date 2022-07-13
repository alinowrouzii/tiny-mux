package modules

import (
	"fmt"
	"net/http"
	"strings"
	"log"
)

var Shit = 50

type radixTree struct {
	root *radixNode
}

type radixNode struct {
	handler http.Handler
	partial string
	childs  map[string]*radixNode
	actualPattern string
}

func (rt *radixTree) insert(urlPattern string, handler http.Handler) {
	partialUrl := partialUrl(urlPattern)
	fmt.Println(partialUrl, len(partialUrl))

	if rt.root == nil {
		fmt.Println("null")
		rootNode := createRadixNode(partialUrl[0], nil)
		rt.root = rootNode
	}

	currNode := rt.root
	for i := 0; i < len(partialUrl)-1; i++ {
		partial := partialUrl[i+1]
		childs := currNode.childs

		if strings.HasPrefix(partial, ":") {
			partial = "*"
		}

		child, childFound := childs[partial]	
		

		if childFound {
			fmt.Println("found")
			currNode = child
		} else{
			fmt.Println("not found")
			newChild := createRadixNode(partial, nil)
			currNode.childs[partial] = newChild
			currNode = newChild
		}
	}

	if currNode.handler == nil {
		currNode.handler = handler
		currNode.actualPattern = urlPattern
	} else {
		log.Fatal(fmt.Sprintf("two url has conflict with each other %s - %s", currNode.actualPattern,urlPattern))
	}
}

func createRadixNode(partial string, handler http.Handler) *radixNode {
	return &radixNode{
		handler: handler,
		partial: partial,
		childs:  make(map[string]*radixNode),
	}
}

func (rt *radixTree) search(urlPattern string) http.Handler {
	partialUrl := partialUrl(urlPattern)
	fmt.Println(partialUrl, len(partialUrl))


	currNode := rt.root
	wildCurrNode := rt.root
	for i := 0; i < len(partialUrl)-1; i++ {
		partial := partialUrl[i+1]
		childs := currNode.childs

		child, childFound := childs[partial]	
		if childFound {
			fmt.Println("found")
			currNode = child
			
		}

		// then looking for wildcard
		wildChild, wildChildFound := childs["*"]

		if wildChildFound {
			fmt.Println("found")
			wildCurrNode = wildChild
		} 
		
		if !childFound && !wildChildFound{
			fmt.Println("partial not found", partial)
			return nil
		}
	}

	if currNode.handler != nil {
		fmt.Println("found pattern", currNode.actualPattern)
		return currNode.handler 
	} 
	
	if wildCurrNode.handler != nil {
		fmt.Println("found wildcardPattern", wildCurrNode.actualPattern)
		return wildCurrNode.handler 
	}
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

	if handler == nil {
		http.NotFound(w, r)
		return
	}

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
	if partialUrl[len(partialUrl)-1]== "" {
		return partialUrl[1 : len(partialUrl)-1]
	}
	return partialUrl[1 : len(partialUrl)]
}
