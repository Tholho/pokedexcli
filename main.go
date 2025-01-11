package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/tholho/pokedexcli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(*config, *pokecache.Cache, string) error
}

type config struct {
	next     string
	previous string
	area     string
	pokedex  map[string]pokemonAPIResponse
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

type locationAPIResponse struct {
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	GameIndex int `json:"game_index"`
	ID        int `json:"id"`
	Location  struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Name  string `json:"name"`
	Names []struct {
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
		Name string `json:"name"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
			MaxChance int `json:"max_chance"`
			Version   struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

type pokemonAPIResponse struct {
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"ability"`
		IsHidden bool `json:"is_hidden"`
		Slot     int  `json:"slot"`
	} `json:"abilities"`
	BaseExperience int `json:"base_experience"`
	Cries          struct {
		Latest string `json:"latest"`
		Legacy string `json:"legacy"`
	} `json:"cries"`
	Forms []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"forms"`
	GameIndices []struct {
		GameIndex int `json:"game_index"`
		Version   struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"version"`
	} `json:"game_indices"`
	Height                 int    `json:"height"`
	HeldItems              []any  `json:"held_items"`
	ID                     int    `json:"id"`
	IsDefault              bool   `json:"is_default"`
	LocationAreaEncounters string `json:"location_area_encounters"`
	Moves                  []struct {
		Move struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"move"`
		VersionGroupDetails []struct {
			LevelLearnedAt  int `json:"level_learned_at"`
			MoveLearnMethod struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"move_learn_method"`
			VersionGroup struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version_group"`
		} `json:"version_group_details"`
	} `json:"moves"`
	Name          string `json:"name"`
	Order         int    `json:"order"`
	PastAbilities []any  `json:"past_abilities"`
	PastTypes     []any  `json:"past_types"`
	Species       struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"species"`
	Sprites struct {
		BackDefault      string `json:"back_default"`
		BackFemale       any    `json:"back_female"`
		BackShiny        string `json:"back_shiny"`
		BackShinyFemale  any    `json:"back_shiny_female"`
		FrontDefault     string `json:"front_default"`
		FrontFemale      any    `json:"front_female"`
		FrontShiny       string `json:"front_shiny"`
		FrontShinyFemale any    `json:"front_shiny_female"`
		Other            struct {
			DreamWorld struct {
				FrontDefault string `json:"front_default"`
				FrontFemale  any    `json:"front_female"`
			} `json:"dream_world"`
			Home struct {
				FrontDefault     string `json:"front_default"`
				FrontFemale      any    `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale any    `json:"front_shiny_female"`
			} `json:"home"`
			OfficialArtwork struct {
				FrontDefault string `json:"front_default"`
				FrontShiny   string `json:"front_shiny"`
			} `json:"official-artwork"`
			Showdown struct {
				BackDefault      string `json:"back_default"`
				BackFemale       any    `json:"back_female"`
				BackShiny        string `json:"back_shiny"`
				BackShinyFemale  any    `json:"back_shiny_female"`
				FrontDefault     string `json:"front_default"`
				FrontFemale      any    `json:"front_female"`
				FrontShiny       string `json:"front_shiny"`
				FrontShinyFemale any    `json:"front_shiny_female"`
			} `json:"showdown"`
		} `json:"other"`
		Versions struct {
			GenerationI struct {
				RedBlue struct {
					BackDefault      string `json:"back_default"`
					BackGray         string `json:"back_gray"`
					BackTransparent  string `json:"back_transparent"`
					FrontDefault     string `json:"front_default"`
					FrontGray        string `json:"front_gray"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"red-blue"`
				Yellow struct {
					BackDefault      string `json:"back_default"`
					BackGray         string `json:"back_gray"`
					BackTransparent  string `json:"back_transparent"`
					FrontDefault     string `json:"front_default"`
					FrontGray        string `json:"front_gray"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"yellow"`
			} `json:"generation-i"`
			GenerationIi struct {
				Crystal struct {
					BackDefault           string `json:"back_default"`
					BackShiny             string `json:"back_shiny"`
					BackShinyTransparent  string `json:"back_shiny_transparent"`
					BackTransparent       string `json:"back_transparent"`
					FrontDefault          string `json:"front_default"`
					FrontShiny            string `json:"front_shiny"`
					FrontShinyTransparent string `json:"front_shiny_transparent"`
					FrontTransparent      string `json:"front_transparent"`
				} `json:"crystal"`
				Gold struct {
					BackDefault      string `json:"back_default"`
					BackShiny        string `json:"back_shiny"`
					FrontDefault     string `json:"front_default"`
					FrontShiny       string `json:"front_shiny"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"gold"`
				Silver struct {
					BackDefault      string `json:"back_default"`
					BackShiny        string `json:"back_shiny"`
					FrontDefault     string `json:"front_default"`
					FrontShiny       string `json:"front_shiny"`
					FrontTransparent string `json:"front_transparent"`
				} `json:"silver"`
			} `json:"generation-ii"`
			GenerationIii struct {
				Emerald struct {
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"emerald"`
				FireredLeafgreen struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"firered-leafgreen"`
				RubySapphire struct {
					BackDefault  string `json:"back_default"`
					BackShiny    string `json:"back_shiny"`
					FrontDefault string `json:"front_default"`
					FrontShiny   string `json:"front_shiny"`
				} `json:"ruby-sapphire"`
			} `json:"generation-iii"`
			GenerationIv struct {
				DiamondPearl struct {
					BackDefault      string `json:"back_default"`
					BackFemale       any    `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  any    `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"diamond-pearl"`
				HeartgoldSoulsilver struct {
					BackDefault      string `json:"back_default"`
					BackFemale       any    `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  any    `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"heartgold-soulsilver"`
				Platinum struct {
					BackDefault      string `json:"back_default"`
					BackFemale       any    `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  any    `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"platinum"`
			} `json:"generation-iv"`
			GenerationV struct {
				BlackWhite struct {
					Animated struct {
						BackDefault      string `json:"back_default"`
						BackFemale       any    `json:"back_female"`
						BackShiny        string `json:"back_shiny"`
						BackShinyFemale  any    `json:"back_shiny_female"`
						FrontDefault     string `json:"front_default"`
						FrontFemale      any    `json:"front_female"`
						FrontShiny       string `json:"front_shiny"`
						FrontShinyFemale any    `json:"front_shiny_female"`
					} `json:"animated"`
					BackDefault      string `json:"back_default"`
					BackFemale       any    `json:"back_female"`
					BackShiny        string `json:"back_shiny"`
					BackShinyFemale  any    `json:"back_shiny_female"`
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"black-white"`
			} `json:"generation-v"`
			GenerationVi struct {
				OmegarubyAlphasapphire struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"omegaruby-alphasapphire"`
				XY struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"x-y"`
			} `json:"generation-vi"`
			GenerationVii struct {
				Icons struct {
					FrontDefault string `json:"front_default"`
					FrontFemale  any    `json:"front_female"`
				} `json:"icons"`
				UltraSunUltraMoon struct {
					FrontDefault     string `json:"front_default"`
					FrontFemale      any    `json:"front_female"`
					FrontShiny       string `json:"front_shiny"`
					FrontShinyFemale any    `json:"front_shiny_female"`
				} `json:"ultra-sun-ultra-moon"`
			} `json:"generation-vii"`
			GenerationViii struct {
				Icons struct {
					FrontDefault string `json:"front_default"`
					FrontFemale  any    `json:"front_female"`
				} `json:"icons"`
			} `json:"generation-viii"`
		} `json:"versions"`
	} `json:"sprites"`
	Stats []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
	Weight int `json:"weight"`
}

var cmdRegistry map[string]cliCommand

func commandExit(config *config, cache *pokecache.Cache, userParam string) error {
	_, err := fmt.Print("Closing the Pokedex... Goodbye!\n")
	if err != nil {
		return err
	}
	os.Exit(0)
	return nil
}

func commandHelp(config *config, cache *pokecache.Cache, userParam string) error {
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

func commandMap(config *config, cache *pokecache.Cache, userParam string) error {
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

func commandMapb(config *config, cache *pokecache.Cache, userParam string) error {
	var jsonData locationAreaAPIResponse
	//fmt.Print(config.previous)
	var data []byte
	if config.previous == "" {
		fmt.Println("you're on the first page")
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

func commandExplore(config *config, cache *pokecache.Cache, location string) error {
	var jsonData locationAPIResponse
	var data []byte
	if location == "" {
		fmt.Println("Please enter a location")
		return fmt.Errorf("no location parameter")
	}
	url := "https://pokeapi.co/api/v2/location-area/" + location
	config.area = location
	if cacheData, ok := cache.Get(url); ok {
		//fmt.Println("using cache")
		data = cacheData
	} else {
		//fmt.Println("fetching data")
		res, err := http.Get(url)
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
	for _, occurrence := range jsonData.PokemonEncounters {
		fmt.Println(occurrence.Pokemon.Name)
	}
	return nil
}

func commandCatch(config *config, cache *pokecache.Cache, pokemon string) error {
	var jsonData pokemonAPIResponse
	var data []byte
	/*
		if config.area == "" {
			fmt.Println("Please explore an area before trying to catch a pokemon")
			return nil
		}
	*/
	url := "https://pokeapi.co/api/v2/pokemon/" + pokemon
	if cacheData, ok := cache.Get(url); ok {
		//fmt.Println("using cache")
		data = cacheData
	} else {
		//fmt.Println("fetching data")
		res, err := http.Get(url)
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
	//suppose max diff is 400
	// observed min diff being 40
	//base chance should be about 35%
	//min chance should be about 1%
	fmt.Print("Throwing a Pokeball at ", pokemon, "...\n")
	pokemonCatchDifficulty := (400 - jsonData.BaseExperience) / 10
	fmt.Println(pokemonCatchDifficulty)
	pokemonCatched := rand.Intn(100) - pokemonCatchDifficulty
	fmt.Println(pokemonCatched)
	if pokemonCatched < 110 {
		fmt.Println(pokemon, "was caught!")
		config.pokedex[pokemon] = jsonData
	} else {
		fmt.Println(pokemon, "escaped!")
	}
	return nil
}

func commandInspect(config *config, cache *pokecache.Cache, pokemon string) error {
	if pokemon == "" {
		fmt.Println("You have not caught that pokemon")
	}
	if value, exists := config.pokedex[pokemon]; exists {
		fmt.Println("Name:", value.Name)
		fmt.Println("Height:", value.Height)
		fmt.Println("Weight:", value.Weight)
		fmt.Println("Stats:")
		for _, val := range value.Stats {
			fmt.Print("	-", val.Stat.Name, ":", val.BaseStat, "\n")
		}
		fmt.Println("Types:")
		for _, val := range value.Types {
			fmt.Print("	-", val.Type.Name, "\n")
		}
	} else {
		fmt.Println("You have not caught that pokemon")
	}
	return nil
}

func commandPokedex(config *config, cache *pokecache.Cache, pokemon string) error {
	fmt.Println("Your pokedex:")
	for _, v := range config.pokedex {
		fmt.Println("-", v.Name)
	}
	return nil
}

func main() {
	var cfgCmd config
	cfgCmd.pokedex = map[string]pokemonAPIResponse{}
	cache := pokecache.NewCache(30 * time.Second)
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
		"explore": {
			name:        "explore",
			description: "Allows the user to see existing pokemon at a given location eg. 'explore location-name' as listed with map command",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Tries to catch a pokemon existing in an area",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "If already caught, displays info about a given pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Displays a list of pokemon in the pokedex",
			callback:    commandPokedex,
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
		var userFirstParameter string
		for registryItem := range cmdRegistry {
			if userCommand == registryItem {
				cmdNotFound = false
				if len(userWords) >= 2 {
					userFirstParameter = userWords[1]
				}
				cmdRegistry[userCommand].callback(&cfgCmd, cache, userFirstParameter)
			}
		}
		if cmdNotFound {
			fmt.Print("Unknown command\n")
		}
		//fmt.Print("Your command was: ", userWords[0], "\n")
	}
}
