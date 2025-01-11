package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/tholho/pokedexcli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, *pokecache.Cache) error
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

func commandExit(config *config, cache *pokecache.Cache) error {
	_, err := fmt.Print("Closing the Pokedex... Goodbye!\n")
	if err != nil {
		return err
	}
	os.Exit(0)
	return nil
}

func commandHelp(config *config, cache *pokecache.Cache) error {
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

func commandMap(config *config, cache *pokecache.Cache) error {
	var jsonData locationAreaAPIResponse
	//fmt.Print(config.previous)
	var data []byte
	if config.previous == "" {
		//	fmt.Print("EMPTY PREVIOUS LOL")
		if cacheData, ok := cache.Get("https://pokeapi.co/api/v2/location-area/"); ok {
			//fmt.Println("using cache")
			data = cacheData
		} else {
			res, err := http.Get("https://pokeapi.co/api/v2/location-area/")
			//fmt.Println("fetching data")
			if err != nil {
				return err
			}
			data, err = io.ReadAll(res.Body)
			if err != nil {
				return err
			}
			cache.Add("https://pokeapi.co/api/v2/location-area/", data)
		}
		err := json.Unmarshal(data, &jsonData)
		if err != nil {
			return err
		}

		config.previous = "https://pokeapi.co/api/v2/location-area/"
		config.next = jsonData.Next
	} else {
		if config.next == "" {
			return fmt.Errorf("either there are no locations left, or an error occured")
		}
		if cacheData, ok := cache.Get(config.next); ok {
			//fmt.Println("using cache")
			data = cacheData
		} else {
			res, err := http.Get(config.next)
			//fmt.Println("fetching data")
			if err != nil {
				return err
			}
			data, err = io.ReadAll(res.Body)
			if err != nil {
				return err
			}
			cache.Add(config.next, data)
		}
		err := json.Unmarshal(data, &jsonData)
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

func commandMapb(config *config, cache *pokecache.Cache) error {
	var jsonData locationAreaAPIResponse
	//fmt.Print(config.previous)
	var data []byte
	if config.previous == "" {
		//fmt.Println("you're on the first page")
		return nil
	} else {
		if cacheData, ok := cache.Get(config.previous); ok {
			//fmt.Println("using cache")
			data = cacheData
		} else {
			//fmt.Println("fetching data")
			res, err := http.Get(config.previous)
			if err != nil {
				return err
			}
			data, err = io.ReadAll(res.Body)
			if err != nil {
				return err
			}
			cache.Add(config.previous, data)
		}
		err := json.Unmarshal(data, &jsonData)
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
	cache := pokecache.NewCache(10 * time.Second)
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
				cmdRegistry[userCommand].callback(&cfgCmd, cache)
			}
		}
		if cmdNotFound {
			fmt.Print("Unknown command\n")
		}
		//fmt.Print("Your command was: ", userWords[0], "\n")
	}
}
