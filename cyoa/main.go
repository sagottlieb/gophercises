package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
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

	var story map[string]chapter
	json.Unmarshal(bytes, &story)

	var someChap chapter
	for _, c := range story {
		someChap = c
		break
	}

	htmlTmpl, err := template.ParseFiles("layout.html")
	if err != nil {
		log.Fatalln(err)
	}

	htmlTmpl.Execute(os.Stdout, someChap)
}
