package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type storyOption struct {
	Text string
	Arc  string
}

type chapter struct {
	Title   string
	Story   []string
	Options []storyOption
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
}
