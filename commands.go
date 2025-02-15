package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"

	"github.com/drewheasman/pokedexcli/internal/pokecache"
)

type pageConfig struct {
	PreviousUrl string
	NextUrl     string
	Cache       pokecache.Cache
	Pokedex     pokedex
}

type cliCommand struct {
	name        string
	description string
	callback    func(*pageConfig, []string) error
}

type pokedex struct {
	caughtPokemon map[string]pokemon
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

	fmt.Println("\n'explore' an area, 'map' for the next page, or 'mapb' the previous page")

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

	fmt.Println("\n'explore' an area, 'map' for the next page, or 'mapb' the previous page")

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

	fmt.Println("\nTry to 'catch' a pokémon...")

	return nil
}

func commandCatch(pageConfig *pageConfig, args []string) error {
	if len(args) == 0 {
		return errors.New("catch requires a pokémon name!")
	}

	pokemon, err := callPokemon(args[0], pageConfig)
	if err != nil {
		return err
	}

	fmt.Printf("throwing a poké ball at %v...\n", pokemon.Name)

	randNum := rand.Float64()
	catchThreshold := float64(pokemon.BaseExperience) / 400.0
	if randNum < catchThreshold {
		fmt.Printf("%v escaped!\n", pokemon.Name)
	} else {
		pageConfig.Pokedex.caughtPokemon[pokemon.Name] = pokemon
		fmt.Printf("%v was caught! Check your 'pokedex'\n", pokemon.Name)
	}

	return nil
}

func commandInspect(pageConfig *pageConfig, args []string) error {
	if len(args) == 0 {
		return errors.New("inspect requires a pokémon name!")
	}

	pokemon, ok := pageConfig.Pokedex.caughtPokemon[args[0]]
	if !ok {
		return fmt.Errorf("%v not found in your pokédex", args[0])
	}

	fmt.Println("Name:", pokemon.Name)
	fmt.Println("Height:", pokemon.Height)
	fmt.Println("Weight:", pokemon.Weight)
	fmt.Println("Stats:")
	for _, s := range pokemon.Stats {
		fmt.Printf(" - %v: %v\n", s.Stat.Name, s.BaseStat+s.Effort)
	}
	fmt.Println("Types:")
	for _, t := range pokemon.Types {
		fmt.Printf(" - %v\n", t.Type.Name)
	}

	return nil
}

func commandPokedex(pageConfig *pageConfig, args []string) error {
	fmt.Println("Your pokédex:")

	if len(pageConfig.Pokedex.caughtPokemon) == 0 {
		fmt.Println("There's nothing here yet! Try catching some pokémon!")
		return nil
	}

	for _, p := range pageConfig.Pokedex.caughtPokemon {
		fmt.Println(" -", p.Name)
	}

	fmt.Println("\nWhy not 'inspect' a pokémon...")

	return nil
}
