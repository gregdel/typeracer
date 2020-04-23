package main

import (
	"encoding/json"
	"image/color"
	"os"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

var (
	colorFg = color.RGBA{R: 0xFD, G: 0xFF, B: 0xFC, A: 255}
	colorBg = color.RGBA{R: 0x01, G: 0x16, B: 0x27, A: 255}
	colorPt = color.RGBA{R: 0xFF, G: 0x9F, B: 0x1C, A: 255}
	colorLn = color.RGBA{R: 0xE7, G: 0x1D, B: 0x36, A: 255}
)

// DataType represents the data type to plot
type dataType string

// Plotable data types
const (
	dataTypeWPM            dataType = "WMP"
	dataTypeAccuracy       dataType = "Accuracy"
	dataTypeWPMPerDday     dataType = "WMP per day"
	dataTypeAccuracyPerDay dataType = "Accuracy per day"
)

func createPlot() error {
	file, err := os.Open("stats.json")
	if err != nil {
		return err
	}
	defer file.Close()

	races := []*Race{}
	if err := json.NewDecoder(file).Decode(&races); err != nil {
		return err
	}

	plots := make([][]*plot.Plot, 2)
	plots[0] = make([]*plot.Plot, 2)
	plots[1] = make([]*plot.Plot, 2)
	for i, d := range []dataType{dataTypeWPM, dataTypeAccuracy} {
		p, err := getPointsPlot(d, races)
		if err != nil {
			return err
		}

		plots[0][i] = p
	}

	for i, d := range []dataType{dataTypeWPMPerDday, dataTypeAccuracyPerDay} {
		p, err := getBarsPlot(d, races)
		if err != nil {
			return err
		}

		plots[1][i] = p
	}

	vgimg.UseBackgroundColor(colorBg)
	img := vgimg.NewWith(
		vgimg.UseWH(vg.Points(910), vg.Points(540)),
		vgimg.UseBackgroundColor(colorBg),
		vgimg.UseDPI(vgimg.DefaultDPI*2),
	)

	dc := draw.New(img)
	rows := 2
	cols := 2
	tiles := draw.Tiles{
		Cols:      rows,
		Rows:      cols,
		PadTop:    vg.Points(10),
		PadRight:  vg.Points(10),
		PadBottom: vg.Points(10),
		PadLeft:   vg.Points(10),
		PadX:      vg.Points(30),
		PadY:      vg.Points(30),
	}

	canvases := plot.Align(plots, tiles, dc)
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			plots[i][j].Draw(canvases[i][j])
		}
	}

	out, err := os.Create("plot.png")
	if err != nil {
		return err
	}
	defer out.Close()

	png := vgimg.PngCanvas{Canvas: img}
	_, err = png.WriteTo(out)
	return err
}
