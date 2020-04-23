package main

import (
	"encoding/json"
	"os"
)

func save() error {
	stats, err := fetch()
	if err != nil {
		return err
	}

	file, err := os.Create("stats.json")
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(stats)
}
