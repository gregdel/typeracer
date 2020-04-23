package main

import (
	"fmt"
	"net/http"

	"github.com/gocolly/colly"
)

const baseURL = "https://data.typeracer.com/pit"

func fetch() ([]*Race, error) {
	races := []*Race{}

	config, err := ConfigFromFile("config.yml")
	if err != nil {
		return nil, err
	}

	c := colly.NewCollector()
	err = c.Post(baseURL+"/login", map[string]string{
		"username": config.Username,
		"password": config.Password,
	})
	if err != nil {
		return nil, err
	}

	// attach callbacks after login
	c.OnResponse(func(r *colly.Response) {
		if r.StatusCode != http.StatusOK {
			fmt.Printf("err: invalid response code: %d", r.StatusCode)
		}
	})

	// Find and visit all links
	c.OnHTML("table.scoresTable", func(e *colly.HTMLElement) {
		e.ForEach("tbody > tr", func(i int, se *colly.HTMLElement) {
			raw := rawRace{}
			se.Unmarshal(&raw)

			if raw.Number == "" {
				return
			}

			race, err := raw.race()
			if err != nil {
				fmt.Printf("err: failed to read race data: %q", err)
			}

			races = append(races, race)
		})
	})

	dataURL := fmt.Sprintf(baseURL+"/race_history?user=%s&n=1000&startDate=&universe=", config.Username)

	c.Visit(dataURL)

	return races, nil
}
