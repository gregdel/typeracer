package main

import (
	"fmt"

	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// XYer wraps the Len and XY methods.
type dataPoints struct {
	x []float64
	y []float64
}

// Len implements the XYer interface
func (d *dataPoints) Len() int {
	return len(d.x)
}

// XY implements the XYer interface
func (d *dataPoints) XY(i int) (x, y float64) {
	return d.x[i], d.y[i]
}

func newDataPoints(d dataType, races []*Race) (*dataPoints, error) {
	data := &dataPoints{
		x: make([]float64, len(races)),
		y: make([]float64, len(races)),
	}
	for i, race := range races {
		data.x[i] = float64(race.Number)
		switch d {
		case dataTypeWPM:
			data.y[i] = float64(race.WPM)
		case dataTypeAccuracy:
			data.y[i] = float64(race.Accuracy)
		default:
			return nil, fmt.Errorf("err: invalid datatype: %s", string(d))
		}
	}

	return data, nil
}

func (d *dataPoints) linearRegrestion() *plotter.Function {
	b, a := stat.LinearRegression(d.x, d.y, nil, false)
	return plotter.NewFunction(func(x float64) float64 { return a*x + b })
}

func getPointsPlot(d dataType, races []*Race) (*plot.Plot, error) {
	data, err := newDataPoints(d, races)
	if err != nil {
		return nil, err
	}

	graph, err := plotter.NewScatter(data)
	if err != nil {
		return nil, err
	}
	graph.Color = colorPt
	graph.Shape = draw.CircleGlyph{}
	graph.GlyphStyle.Radius = vg.Points(1)

	regression := data.linearRegrestion()
	regression.Color = colorLn
	regression.Width = vg.Points(2)

	p, err := plot.New()
	if err != nil {
		return nil, err
	}
	p.Title.Text = string(d)
	p.Title.Color = colorFg
	p.BackgroundColor = colorBg

	p.X.Tick.Color = colorFg
	p.X.Tick.Label.Color = colorFg
	p.X.Color = colorFg

	p.Y.Tick.Color = colorFg
	p.Y.Tick.Label.Color = colorFg
	p.Y.Color = colorFg

	p.Add(graph, regression)
	return p, nil
}
