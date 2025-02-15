package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
)

const baseUrl = "https://pokeapi.co/api/v2/"

type locationAreas struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type locationAreasCallType string

const (
	Back    locationAreasCallType = "back"
	Forward locationAreasCallType = "forward"
)

func callLocationAreas(callType locationAreasCallType, pageConfig *pageConfig) ([]string, error) {
	var areaNames []string

	locationAreaUrlString := baseUrl + "location-area"
	if callType == Back {
		if len(pageConfig.PreviousUrl) == 0 {
			return areaNames, fmt.Errorf("locationAreasCallType Back, but no Back url")
		}
		locationAreaUrlString = pageConfig.PreviousUrl
	}
	if callType == Forward {
		if len(pageConfig.NextUrl) > 0 {
			locationAreaUrlString = pageConfig.NextUrl
		}
	}

	locationAreaUrl, err := url.Parse(locationAreaUrlString)
	if err != nil {
		return areaNames, fmt.Errorf("locationAreaUrl invalid")
	}

	data, err := cacheOrCall(&pageConfig.Cache, locationAreaUrl)
	if err != nil {
		return areaNames, err
	}

	var respJson locationAreas
	if err := json.Unmarshal(data, &respJson); err != nil {
		return areaNames, fmt.Errorf("error unmarshalling location-area response")
	}

	for _, result := range respJson.Results {
		areaNames = append(areaNames, result.Name)
	}
	pageConfig.PreviousUrl = respJson.Previous
	pageConfig.NextUrl = respJson.Next

	return areaNames, nil
}

type locationArea struct {
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

func callLocationAreaId(id string, pageConfig *pageConfig) ([]string, error) {
	var pokemonNames []string

	locationAreaUrlString := baseUrl + "location-area" + "/" + id

	locationAreaUrl, err := url.Parse(locationAreaUrlString)
	if err != nil {
		return pokemonNames, fmt.Errorf("locationAreaUrl invalid")
	}

	data, err := cacheOrCall(&pageConfig.Cache, locationAreaUrl)
	if err != nil {
		return pokemonNames, err
	}

	var respJson locationArea
	if err := json.Unmarshal(data, &respJson); err != nil {
		return pokemonNames, fmt.Errorf("error unmarshalling location-area id response")
	}

	for _, encounter := range respJson.PokemonEncounters {
		pokemonNames = append(pokemonNames, encounter.Pokemon.Name)
	}

	return pokemonNames, nil
}

type pokemon struct {
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
	Sprites map[string]interface{} `json:"sprites"`
	Stats   []struct {
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

func callPokemon(name string, pageConfig *pageConfig) (pokemon, error) {
	var pokemonster pokemon

	if len(name) == 0 {
		return pokemonster, errors.New("pokemon not specified")
	}

	pokemonUrlString := baseUrl + "pokemon" + "/" + name

	pokemonUrl, err := url.Parse(pokemonUrlString)
	if err != nil {
		return pokemonster, fmt.Errorf("pokemonUrl invalid")
	}

	data, err := cacheOrCall(&pageConfig.Cache, pokemonUrl)
	if err != nil {
		return pokemonster, err
	}

	if err := json.Unmarshal(data, &pokemonster); err != nil {
		return pokemonster, fmt.Errorf("error unmarshalling pokemon response")
	}

	return pokemonster, nil
}
