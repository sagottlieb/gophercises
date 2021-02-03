package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

	// For simplicity, all stories will have a story arc named "intro" that is where the story starts.
	err := tmpl.Execute(w, h.story["intro"])
	if err != nil {
		log.Fatalln(err)
	}
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
