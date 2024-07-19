package main

import "fmt"
import "sync"
import "time"
import "os"
import "flag"
import "encoding/json"

var Config Configuration

var apiKey = flag.String("api-key", "", "Enter the api key for the Merriam-Webster api")

func main() {

	// Parse command line flags
	flag.Parse()
	
	// Get the path to the user's home directory
	homePath, ok := os.LookupEnv("HOME")
	if !ok {
		fmt.Println("Couldn't find a path to the home directory")
		return
	}

	file, err := os.OpenFile(homePath + "/.wotdconfig", os.O_RDWR | os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error opening file")
		fmt.Println(err)
		return
	}

	bytes := BytesFromReader(file, 512, -1)
	configString := BytesToString(bytes)

	if configString == "" {
		if *apiKey == "" {
			fmt.Println("Need an api key")
			return
		}

		Config.Key = *apiKey
		bytes, err = json.Marshal(Config)
		if err != nil {
			fmt.Println("Error marshaling to json")
			return
		}
		
		_, err = file.Write(bytes)
		if err != nil {
			fmt.Println("Error writing to file")
			return
		}
	} else {
		if *apiKey != "" {
			Config.Key = *apiKey
			bytes, err = json.Marshal(Config)
			if err != nil {
				fmt.Println("Error marshaling json")
				fmt.Println(err)
				return
			}

			
			err = file.Truncate(0)
			file.Seek(0, 0)
			if err != nil {
				fmt.Println("Error truncating file")
				return
			}

			_, err = file.Write(bytes)
			if err != nil {
				fmt.Println("Error writing to file")
				return
			}
		} else {
			var loadedConfig Configuration
			err = json.Unmarshal(bytes, &loadedConfig)
			if err != nil {
				fmt.Println("Error unmarshaling config")
				fmt.Println(err)
				return
			}

			Config.Key = loadedConfig.Key
		}
	}	

	file.Close()
	

	// Variables to store the definitions and error from GetDefinition
	var words *[]WordData
	// Define and lock a mutex
	var mut sync.Mutex
	mut.Lock() // lock the mutex outside the goroutine because of interleaving errors
	// Load the definition on a separate thread
	go func() {
		// unlock mutex when goroutine finishes
		defer mut.Unlock()
		words, err = GetDefinition()
	}()

	// Print out a loading animation until the mutex is unlocked
	count := 0
	fmt.Print("Loading: |")
	for !mut.TryLock() {
		n := count % 4
		// \033 is the escape sequence for CSI (Control Sequence Introducer) sequences
		// https://en.wikipedia.org/wiki/ANSI_escape_code
		// \033[1D moves the cursor back 1 cell
		fmt.Print("\033[1D")
		switch n {
		case 0:
			fmt.Print("/")
		case 1:
			fmt.Print("-")
		case 2:
			fmt.Print("\\")
		case 3:
			fmt.Print("|")
		}
		
		count++
		time.Sleep(150 * time.Millisecond)
	}

	// move the cursor back 10 cells to delete the loading stuff
	fmt.Print("\033[10D")

	if err != nil {
		fmt.Println("There has been an error.")
		fmt.Println(err)
		return
	}

	for _, v := range *words {
		fmt.Println(v)
	}
	
}
