package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/drewheasman/pokedexcli/internal/pokecache"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	paging := pageConfig{
		Cache: *pokecache.NewCache(60_000 * time.Millisecond),
	}
	paging.Pokedex.caughtPokemon = make(map[string]pokemon)

	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			panic("fatal while reading input!")
		}

		inputFields := strings.Fields(scanner.Text())
		if len(inputFields) == 0 {
			continue
		}

		command := commands()[inputFields[0]]
		fmt.Println()
		if command.callback == nil {
			println("Unknown command")
		} else {
			if err := command.callback(&paging, inputFields[1:]); err != nil {
				fmt.Println("error:", err)
			}
		}
		fmt.Println()
	}
}

func cleanInput(text string) []string {
	words := strings.Fields(text)
	var lowerWords []string
	for _, w := range words {
		lowerWords = append(lowerWords, strings.ToLower(w))
	}
	return lowerWords
}

func commands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandUsage,
		},
		"map": {
			name:        "map",
			description: "Show the next location areas",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Show the previous location areas",
			callback:    commandMapBack,
		},
		"explore": {
			name:        "explore",
			description: "Explore a location area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Catch a pokémon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a pokémon you have caught",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List the pokémon in your pokédex",
			callback:    commandPokedex,
		},
	}
}
