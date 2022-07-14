package tinymux

import (
	"net/http"
	"testing"
)

func BenchmarkOne(b *testing.B) {

	tm := NewTinyMux()
	tm.Handle("GET", "/foo/:bar", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	req, err := http.NewRequest("GET", "/foo/baar", nil)
	if err != nil {
		b.Fatal(err)
	}
	for n := 0; n < b.N; n++ {
		tm.ServeHTTP(nil, req)
	}
}

func BenchmarkTwo(b *testing.B) {

	tm := NewTinyMux()
	tm.Handle("GET", "/foo/:bar/:baz/:qux/:fred/:thud/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	req, err := http.NewRequest("GET", "/foo/baaaar/baaaaz/quuuux/freeeed/thuuuud/", nil)
	if err != nil {
		b.Fatal(err)
	}
	for n := 0; n < b.N; n++ {
		tm.ServeHTTP(nil, req)
	}
}
