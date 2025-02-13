package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/drewheasman/pokedexcli/internal/pokecache"
)

type pageConfig struct {
	PreviousUrl string
	NextUrl     string
	Cache       pokecache.Cache
}

type cliCommand struct {
	name        string
	description string
	callback    func(*pageConfig, []string) error
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	paging := pageConfig{
		Cache: *pokecache.NewCache(60_000 * time.Millisecond),
	}

	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			panic("fatal while reading input!")
		}

		inputFields := strings.Fields(scanner.Text())
		command := commands()[inputFields[0]]
		fmt.Println()
		if command.callback == nil {
			println("Unknown command")
		} else {
			if err := command.callback(&paging, inputFields[1:]); err != nil {
				fmt.Println("error running command:", err)
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
	}
}

func commandExit(pageConfig *pageConfig, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandUsage(pageConfig *pageConfig, args []string) error {
	fmt.Printf("Welcome to the Pokedex!\n\nUsage:\n\n")
	for _, c := range commands() {
		fmt.Printf("%v: %v\n", c.name, c.description)
	}

	return nil
}

func commandMap(pageConfig *pageConfig, args []string) error {
	areas, err := callLocationAreas(Forward, pageConfig)
	if err != nil {
		return err
	}

	for _, a := range areas {
		fmt.Println(a)
	}

	return nil
}

func commandMapBack(pageConfig *pageConfig, args []string) error {
	if len(pageConfig.PreviousUrl) == 0 {
		fmt.Println("you're on the first page")
		return nil
	}

	areas, err := callLocationAreas(Back, pageConfig)
	if err != nil {
		return err
	}

	for _, a := range areas {
		fmt.Println(a)
	}

	return nil
}

func commandExplore(pageConfig *pageConfig, args []string) error {
	if len(args) == 0 {
		return errors.New("explore requires an area name!")
	}

	pokemonNames, err := callLocationAreaId(args[0], pageConfig)
	if err != nil {
		return err
	}

	for _, n := range pokemonNames {
		fmt.Println(n)
	}

	return nil
}
