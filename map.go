package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const baseUrl = "https://pokeapi.co/api/v2/"

type pageConfig struct {
	PreviousUrl string
	NextUrl     string
}

type locationArea struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type callType string

const (
	Back    callType = "back"
	Forward callType = "forward"
)

func callLocationArea(callType callType, pageConfig *pageConfig) ([]string, error) {
	var areaNames []string

	locationAreaUrl := baseUrl + "location-area"
	if callType == Back {
		if len(pageConfig.PreviousUrl) == 0 {
			return areaNames, fmt.Errorf("location-area callType Back, but no Back url")
		}
		locationAreaUrl = pageConfig.PreviousUrl
	}
	if callType == Forward {
		if len(pageConfig.NextUrl) > 0 {
			locationAreaUrl = pageConfig.NextUrl
		}
	}

	resp, err := http.Get(locationAreaUrl)
	if err != nil {
		return areaNames, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return areaNames, fmt.Errorf("location-area response was not OK")
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return areaNames, fmt.Errorf("error reading location-area response")
	}

	var respJson locationArea
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
