package main

import "fmt"

type Configuration struct {
	// Key is the api key for the Merriam-Webster api
	Key string `json:"key"`
}

type WordData struct {
	// Word is the word being defined
	Word HeadwordInformation `json:"hwi"`
	// ShortDefinitions is a list of abbreviated definitions for the word
	ShortDefinitions []string `json:"shortdef"`
}

func (w WordData) String() string {
	result := "The word of the day is: " + w.Word.String()
	for i, str := range w.ShortDefinitions {
		result += fmt.Sprintf("\n%d) %s", i + 1, str)
	}
	return result
}

type HeadwordInformation struct {
	// Headword is the actual word
	Headword string `json:"hw"`
}

func (h HeadwordInformation) String() string {
	return h.Headword
}
