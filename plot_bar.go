package main

import (
	"fmt"
	"math"
	"sort"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

type dataBar struct {
	// Days in unix timestamp format
	x []string
	// Avg per day
	y []float64
}

// Len implements the Valuer interface
func (d *dataBar) Len() int {
	return len(d.x)
}

// Value implements the Valuer interface
func (d *dataBar) Value(i int) float64 {
	return d.y[i]
}

func (d *dataBar) MinY() float64 {
	minY := d.y[0]
	for _, y := range d.y {
		if y < minY {
			minY = y
		}
	}
	return minY
}

func newDataBar(d dataType, races []*Race) (*dataBar, error) {
	data := &dataBar{}

	dates := []int64{}
	type bar struct {
		date   time.Time
		value  float64
		weight float64
	}
	avgPerDay := map[int64]*bar{}

	for _, race := range races {
		date := race.Date.Unix()
		avg, ok := avgPerDay[date]
		if !ok {
			dates = append(dates, date)
			avg = &bar{date: race.Date, value: 0, weight: 0}
			avgPerDay[date] = avg
		}

		var value float64
		switch d {
		case dataTypeWPMPerDday:
			value = float64(race.WPM)
		case dataTypeAccuracyPerDay:
			value = float64(race.Accuracy)
		default:
			return nil, fmt.Errorf("err: invalid datatype: %s", string(d))
		}

		avgPerDay[date].value = (avg.value*avg.weight + value) / (avg.weight + 1)
		avgPerDay[date].weight = avgPerDay[date].weight + 1

	}

	sort.Slice(dates, func(i, j int) bool { return dates[j]-dates[i] > 0 })

	days := len(dates)
	if days > 10 {
		dates = dates[days-10:]
	} else {
		dates = dates[:days]
	}

	data.x = make([]string, len(dates))
	data.y = make([]float64, len(dates))
	for i, date := range dates {
		data.x[i] = fmt.Sprintf(
			"%s\n(%.0f races)",
			avgPerDay[date].date.Format("Jan 2"),
			avgPerDay[date].weight)
		data.y[i] = avgPerDay[date].value
	}

	return data, nil
}

func getBarsPlot(d dataType, races []*Race) (*plot.Plot, error) {
	data, err := newDataBar(d, races)
	if err != nil {
		return nil, err
	}

	bars, err := plotter.NewBarChart(data, vg.Points(30))
	if err != nil {
		return nil, err
	}
	bars.Color = colorPt
	bars.LineStyle.Color = colorPt

	p, err := plot.New()
	if err != nil {
		return nil, err
	}
	p.Add(bars)
	p.Title.Text = string(d)
	p.Title.Color = colorFg
	p.NominalX(data.x...)
	p.BackgroundColor = colorBg

	p.X.Tick.Color = colorFg
	p.X.Tick.Label.Color = colorFg
	p.X.Tick.Label.Rotation = math.Pi / 4
	p.X.Tick.Label.XAlign = draw.XRight
	p.X.Tick.Label.YAlign = draw.YTop
	p.X.Color = colorFg

	p.Y.Tick.Color = colorFg
	p.Y.Tick.Label.Color = colorFg
	p.Y.Color = colorFg
	p.Y.Min = data.MinY() - 1

	return p, nil
}
