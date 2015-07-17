package tps

import (
	"errors"
)

// Grid holds all the page and grid specification required for the Report to
// create new pages and place content.
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

// CalculateColumns calculates the remaining column related specs based on
// ColumnCount and GutterWidth.
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
