package main

import (
	"fmt"
	"github.com/sagottlieb/gophercises/url-shortener/urlshort"
	"net/http"
)

func main() {
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yaml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
- path: /georgie
  url: https://www.google.com/imgres?imgurl=https%3A%2F%2Fi.pinimg.com%2Foriginals%2F28%2F8b%2F62%2F288b6232a672581fa7199e2cf196db9e.jpg&imgrefurl=https%3A%2F%2Fwww.pinterest.com%2Fpin%2F189925309264432634%2F&docid=Fo5C8zd6XKY-oM&tbnid=0eqR4EBJYrtj-M%3A&vet=1&w=300&h=232&bih=674&biw=1332&ved=0ahUKEwjR76KVg4zbAhUCwVkKHYAeB_QQMwg8KAMwAw&iact=c&ictx=1
`

	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting the server on :8081")
	http.ListenAndServe(":8081", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
