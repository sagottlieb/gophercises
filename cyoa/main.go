package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
)

type storyOption struct {
	Text        string
	ChapterName string `json:"arc"`
}

type chapter struct {
	Title      string
	Paragraphs []string `json:"story"`
	Options    []storyOption
}

type cyoaStory map[string]chapter

func NewHandler(story cyoaStory) http.Handler {
	return handler{story: story}
}

type handler struct {
	story cyoaStory
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("layout.html"))

	path := strings.TrimLeft(strings.TrimSpace(r.URL.Path), "/")

	// every JSON file will have a key with the value `intro` and this is where your story should start.
	if path == "" {
		path = "intro"
	}
	if chapter, exists := h.story[path]; exists {
		err := tmpl.Execute(w, chapter)
		if err != nil {
			log.Printf("%v\n", err)
			http.Error(w, "An adventure indeed...", http.StatusInternalServerError)
		}
		return
	}
	http.Error(w, "Chapter not found", http.StatusNotFound)
}

func main() {
	jsonFile, err := os.Open("gopher.json")
	if err != nil {
		log.Fatalln(err)
	}
	defer jsonFile.Close()

	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatalln(err)
	}

	var story cyoaStory
	json.Unmarshal(bytes, &story)

	h := NewHandler(story)

	fmt.Println("Starting the server on :8081")

	log.Fatalln(http.ListenAndServe(":8081", h))
}
