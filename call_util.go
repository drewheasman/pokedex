package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/drewheasman/pokedexcli/internal/pokecache"
)

func cacheOrCall(cache *pokecache.Cache, url *url.URL) ([]byte, error) {
	var data []byte
	if d, ok := cache.Get(url.String()); ok {
		data = d
	} else {
		resp, err := http.Get(url.String())
		if err != nil {
			return data, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return data, fmt.Errorf("location-area id response was not OK")
		}

		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return data, fmt.Errorf("error reading location-area id response")
		}

		cache.Add(url.String(), data)
	}
	return data, nil
}
