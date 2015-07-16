package tps

import (
	"errors"
)

type Grid struct {
	ColumnCount int
	ColumnWidth float64
	GutterCount int
	GutterWidth float64
	LineHeight  float64
	Margin      float64
	Orientation string
	Size        string
	Units       string
	PageWidth   float64
	PageHeight  float64
}

func (g *Grid) CalculateColumns() error {
	if g.ColumnCount == 0 || g.GutterWidth == 0 {
		return errors.New("Incomplete data to calculate grid columns")
	}
	g.GutterCount = g.ColumnCount - 1
	width := g.PageWidth
	width -= g.Margin * 2
	width -= float64(g.GutterCount) * g.GutterWidth
	g.ColumnWidth = width / float64(g.ColumnCount)
	return nil
}
