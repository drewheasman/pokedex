package main

import (
	"bufio"
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
	callback    func(*pageConfig) error
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	paging := pageConfig{
		Cache: *pokecache.NewCache(10_000 * time.Millisecond),
	}

	for {
		fmt.Print("Pokedex > ")
		if !scanner.Scan() {
			panic("fatal while reading input!")
		}

		command := commands()[scanner.Text()]
		fmt.Println()
		if command.callback == nil {
			println("Unknown command")
		} else {
			command.callback(&paging)
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
	}
}

func commandExit(pageConfig *pageConfig) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandUsage(pageConfig *pageConfig) error {
	fmt.Printf("Welcome to the Pokedex!\nUsage:\n")
	for _, c := range commands() {
		fmt.Printf("%v: %v\n", c.name, c.description)
	}

	return nil
}

func commandMap(pageConfig *pageConfig) error {
	areas, err := callLocationArea(Forward, pageConfig)
	if err != nil {
		return err
	}

	for _, a := range areas {
		fmt.Println(a)
	}

	return nil
}

func commandMapBack(pageConfig *pageConfig) error {
	if len(pageConfig.PreviousUrl) == 0 {
		fmt.Println("you're on the first page")
		return nil
	}

	areas, err := callLocationArea(Back, pageConfig)
	if err != nil {
		return err
	}

	for _, a := range areas {
		fmt.Println(a)
	}

	return nil
}
