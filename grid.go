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
	Orientation int
	PageWidth   float64
	PageHeight  float64
	PageSize    int
	Unit        int
}

// Point is the X, Y coordinates in the Grid.Unit system relative to the PDF's
// coordinate system converted from tps' spreadsheet-like coordinates.
type Point struct {
	X, Y float64
}

// Cell is the width and height in the Grid.Unit system relative to the PDF's
// measurement system converted from tps' speadsheet-like sizing.
type Cell struct {
	Width, Height float64
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

// GetCell returns a Cell struct for use in lower level Fpdf calls
func (g *Grid) GetCell(block Block) Cell {
	cell := Cell{}

	cell.Width = g.ColumnWidth * float64(block.Width)
	cell.Width += g.GutterWidth * float64(block.Width-1)

	cell.Height = g.LineHeight * float64(block.Height)

	return cell
}

// GetPoint returns a Point struct for use in lower level Fpdf calls
func (g *Grid) GetPoint(x, y int) Point {
	point := Point{}

	point.X = g.Margin
	point.X += g.ColumnWidth * float64(x-1)
	point.X += g.GutterWidth * float64(x-1)

	point.Y = g.Margin
	point.Y += g.LineHeight * float64(y-1)

	return point
}

func (g *Grid) convertOrientation() string {
	return orientation[g.Orientation]
}

func (g *Grid) convertPageSize() string {
	return pageSize[g.PageSize]
}

func (g *Grid) convertUnit() string {
	return unit[g.Unit]
}
