package main

import "time"

// Race represents a race
type Race struct {
	Number     int       `json:"number"`
	WPM        int       `json:"wpm"`
	Accuracy   float64   `json:"accuracy"`
	Points     int       `json:"points"`
	Rank       int       `json:"rank"`
	RacerCount int       `json:"racer_count"`
	Date       time.Time `json:"date"`
}
