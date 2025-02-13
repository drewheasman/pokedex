package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

	locationAreaUrl := baseUrl + "location-area"
	if callType == Back {
		if len(pageConfig.PreviousUrl) == 0 {
			return areaNames, fmt.Errorf("locationAreasCallType Back, but no Back url")
		}
		locationAreaUrl = pageConfig.PreviousUrl
	}
	if callType == Forward {
		if len(pageConfig.NextUrl) > 0 {
			locationAreaUrl = pageConfig.NextUrl
		}
	}

	var data []byte
	if d, ok := pageConfig.Cache.Get(locationAreaUrl); ok {
		data = d
	} else {
		resp, err := http.Get(locationAreaUrl)
		if err != nil {
			return areaNames, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return areaNames, fmt.Errorf("location-area response was not OK")
		}

		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return areaNames, fmt.Errorf("error reading location-area response")
		}

		pageConfig.Cache.Add(locationAreaUrl, data)
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

	locationAreaUrl := baseUrl + "location-area" + "/" + id

	var data []byte
	if d, ok := pageConfig.Cache.Get(locationAreaUrl); ok {
		data = d
	} else {
		resp, err := http.Get(locationAreaUrl)
		if err != nil {
			return pokemonNames, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return pokemonNames, fmt.Errorf("location-area id response was not OK")
		}

		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return pokemonNames, fmt.Errorf("error reading location-area id response")
		}

		pageConfig.Cache.Add(locationAreaUrl, data)
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
