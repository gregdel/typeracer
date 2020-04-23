package main

import (
	"strconv"
	"strings"
	"time"
)

type rawRace struct {
	Number   string `selector:"td:nth-child(1) > a"`
	WPM      string `selector:"td:nth-child(2)"`
	Accuracy string `selector:"td:nth-child(3)"`
	Points   string `selector:"td:nth-child(4)"`
	Rank     string `selector:"td:nth-child(5)"`
	Date     string `selector:"td:nth-child(6)"`
}

func (rr *rawRace) race() (*Race, error) {
	race := &Race{}

	// Date
	if rr.Date == "today" {
		now := time.Now()
		race.Date = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	} else {
		date, err := time.Parse("January 02, 2006", rr.Date)
		if err != nil {
			return nil, err
		}
		race.Date = date
	}

	// Accuracy
	accuracy, err := strconv.ParseFloat(strings.Trim(rr.Accuracy, "%"), 64)
	if err != nil {
		return nil, err
	}
	race.Accuracy = accuracy

	// Rank
	var rank, racerCount string
	parts := strings.Split(rr.Rank, "/")
	if len(parts) == 2 {
		rank = parts[0]
		racerCount = parts[1]
	}

	for _, e := range []struct {
		s string
		i *int
	}{
		{s: rr.Number, i: &race.Number},
		{s: strings.Trim(rr.WPM, " WPM"), i: &race.WPM},
		{s: rr.Points, i: &race.Points},
		{s: rank, i: &race.Rank},
		{s: racerCount, i: &race.RacerCount},
	} {
		value, err := strconv.Atoi(e.s)
		if err != nil {
			return nil, err
		}

		*e.i = value
	}

	return race, nil
}
