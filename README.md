### TinyMux, a fast HTTP router for go web server based on radixTree

To download the module:
```
     go get -u github.com/alinowrouzii/tiny-mux
```

#### Usage:
```
    tm := tinymux.NewTinyMux()

    tm.Handle("GET", "/foo", http.HandlerFunc(fooHandler))
    tm.POST("/foo", http.HandlerFunc(fooHandler))

    http.ListenAndServe(":8000", tm)
```

#### Middleware support:
```
    // functions that is wrapper for tm, accepts http.Handler and also returns http.Handler
    tm.Use(middleware1, middleware2, middlewareN)
```
It is important where tm.Use is located.
For-instance two examples below is not the same
```
    tm.Use(middleware1)
    tm.Handle("GET", "/foo", http.HandlerFunc(fooHandler))
    tm.Use(middleware2)
```
```
    tm.Use(middleware1)
    tm.Use(middleware2)
    tm.Handle("GET", "/foo", http.HandlerFunc(fooHandler))
```
In the first example middleware2 does not apply to the fooHandler. But in the second one applies.

#### Named route parameters support:
```
    tm.Handle("GET", "/foo/:bar/:baz", http.HandlerFunc(fooHandler))    
```
To access parameters inside handler:
```
func fooHandler(w http.ResponseWriter, r *http.Request) {
    // suppose we are in "/foo/:bar/:baz" handler
    // and the requested url is "/foo/baar/baaz"

    params := tinymux.Values(*r)
    fmt.Println(params["bar"]) // expected baar
    fmt.Println(params["baz"]) // expected baaz
}
```
#### Apply middleware to a specific handler:
```
    handler := tinymux.ChainMiddlewares(http.HandlerFunc(fooHandler), middleware1, middleware2, middlewareN)
    tm.Handle("/foo/:bar/:baz", handler)    

```
