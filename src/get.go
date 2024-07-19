package main

import "fmt"
import "net/http"
import "strings"
import "encoding/json"
import "errors"

const wotdPage = "https://www.merriam-webster.com/word-of-the-day"
const query = "https://www.dictionaryapi.com/api/v3/references/collegiate/json/{{word}}?key={{api-key}}"

const start = "<title>Word of the Day: "
const end = " | Merriam-Webster</title>"

// GetWord scrapes the word of the day from the Merriam-Webster Word of the Day
// web page and returns it as a string
func GetWord() (string, error) {
	resp, err := http.Get(wotdPage)
	if err != nil {
		fmt.Println("Error getting word of the day: ", err)
		return "", err
	}

	p := BytesFromReader(resp.Body, 2048, 1)
	resp.Body.Close()

	result := BytesToString(p)

	wordOfTheDay := strings.SplitAfter(result, start)[1]
	wordOfTheDay = strings.Split(wordOfTheDay, end)[0]

	return wordOfTheDay, nil
}

// GetDefintion queries the Merriam-Webster api and returns a slice containing
// WordData objects
func GetDefinition() (*[]WordData, error)  {
	word, err := GetWord()
	if err != nil {
		fmt.Println("Error getting the word of the day: ", err)
		return nil, err
	}

	specificQuery := strings.Replace(query, "{{word}}", strings.ToLower(word), 1)
	specificQuery = strings.Replace(specificQuery, "{{api-key}}", Config.Key, 1)

	resp, err := http.Get(specificQuery)
	if err != nil {
		fmt.Println("Error getting json from api: ", err)
		return nil, err
	}

	whole := BytesFromReader(resp.Body, 1024, -1)
	resp.Body.Close()

	if BytesToString(whole) == "Invalid API key. Not subscribed for this reference." {
		fmt.Println("The API key you provided is not valid. Please provide another API key.")
		return nil, errors.New("invalid api key")
	}

	var def []WordData
	err = json.Unmarshal(whole, &def)
	if err != nil {
		fmt.Println("Error unmarshaling json: ", err)
		return nil, err
	}

	return &def, nil
	
}
