package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

type config struct {
	next     string
	previous string
}

type locationAreaAPIResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

var cmdRegistry map[string]cliCommand

func cleanInput(text string) []string {
	result := strings.Fields(text)
	fmt.Print(result)
	return result
}

func commandExit(config *config) error {
	_, err := fmt.Print("Closing the Pokedex... Goodbye!\n")
	if err != nil {
		return err
	}
	os.Exit(0)
	return nil
}

func commandHelp(config *config) error {
	cmdDescriptions := ""
	for item := range cmdRegistry {
		cmdDescriptions = cmdDescriptions + "\n" + cmdRegistry[item].name + ": " + cmdRegistry[item].description
	}
	_, err := fmt.Println("Welcome to the Pokedex!", "Usage:", "", cmdDescriptions)
	if err != nil {
		return err
	}
	return nil
}

func commandMap(config *config) error {
	var jsonData locationAreaAPIResponse
	//fmt.Print(config.previous)
	if config.previous == "" {
		//	fmt.Print("EMPTY PREVIOUS LOL")
		res, err := http.Get("https://pokeapi.co/api/v2/location-area/")
		if err != nil {
			return err
		}
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		err = json.Unmarshal(data, &jsonData)
		if err != nil {
			return err
		}
		config.previous = "https://pokeapi.co/api/v2/location-area/"
		config.next = jsonData.Next
	} else {
		if config.next == "" {
			return fmt.Errorf("either there are no locations left, or an error occured")
		}
		res, err := http.Get(config.next)
		if err != nil {
			return err
		}
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		err = json.Unmarshal(data, &jsonData)
		if err != nil {
			return err
		}

		config.previous = config.next
		config.next = jsonData.Next
	}
	for _, location := range jsonData.Results {
		fmt.Println(location.Name)
	}
	//var data string
	//decoder := json.NewDecoder(res.Body)
	//err = decoder.Decode(&data)
	//fmt.Println(jsonData["previous"])
	//fmt.Println(jsonData["next"])

	fmt.Println(jsonData.Previous)
	fmt.Println(jsonData.Next)
	return nil
}

func commandMapb(config *config) error {
	var jsonData locationAreaAPIResponse
	//fmt.Print(config.previous)
	if config.previous == "" {
		fmt.Println("you're on the first page")
		return nil
	} else {
		res, err := http.Get(config.previous)
		if err != nil {
			return err
		}
		data, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		err = json.Unmarshal(data, &jsonData)
		if err != nil {
			return err
		}
		config.next = config.previous
		config.previous = jsonData.Previous

	}
	for _, location := range jsonData.Results {
		fmt.Println(location.Name)

	}
	fmt.Println(jsonData.Previous)
	fmt.Println(jsonData.Next)
	return nil
}

func main() {
	var cfgCmd config
	cmdRegistry = map[string]cliCommand{
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Lists 20 Poke-Locations... and the next 20 ones for each subsequent commands",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Lists previous map results, if exist",
			callback:    commandMapb,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
	}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex >")
		if !scanner.Scan() {
			fmt.Println("No more input. Exiting.")
			break
		}
		userCurrentInput := scanner.Text()
		loweredUserInput := strings.ToLower(userCurrentInput)
		userWords := strings.Fields(loweredUserInput)
		if len(userWords) == 0 {
			fmt.Println("Please enter a valid command.")
			continue
		}
		userCommand := userWords[0]
		cmdNotFound := true
		for registryItem := range cmdRegistry {
			if userCommand == registryItem {
				cmdNotFound = false
				cmdRegistry[userCommand].callback(&cfgCmd)
			}
		}
		if cmdNotFound {
			fmt.Print("Unknown command\n")
		}
		//fmt.Print("Your command was: ", userWords[0], "\n")
	}
}
