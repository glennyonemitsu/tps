package tps

import (
	"testing"
)

func newGrid() Grid {
	g := Grid{
		ColumnCount: 12,
		GutterWidth: 12.0,
		LineHeight:  12.0,
		Margin:      36.0,
		Orientation: OrientationPortrait,
		PageWidth:   612.0,
		PageHeight:  792.0,
		PageSize:    PageSizeLetter,
		Unit:        UnitPt,
	}
	g.CalculateColumns()
	// ColumnWidth == 34.0
	return g
}

func TestCalculateColumns(t *testing.T) {
	g := Grid{
		ColumnCount: 6,
		Margin:      20.0,
		PageWidth:   240.0,
	}

	err := g.CalculateColumns()
	if err == nil {
		t.Error(err)
	}

	g.GutterWidth = 10.0
	err = g.CalculateColumns()
	if err != nil {
		t.Error("Grid calculate failed when required params were provided.")
	}

	if g.GutterCount != 5 {
		t.Error("Grid did not calculate proper GutterCount.")
	}

	if g.ColumnWidth != 25.0 {
		t.Error("Grid did not calculate proper ColumnWidth.")
	}

}

func TestConvertOrientation(t *testing.T) {
	g := Grid{
		Orientation: OrientationPortrait,
	}
	if g.convertOrientation() != "Portrait" {
		t.Error("Grid did not interpret orientation constant correctly.")
	}
	g.Orientation = -1
	if g.convertOrientation() != "" {
		t.Error("Grid did not return empty string for improper orientation constant.")
	}
}

func TestConvertUnit(t *testing.T) {
	g := Grid{
		Unit: UnitPt,
	}
	if g.convertUnit() != "pt" {
		t.Error("Grid did not interpret Unit constant correctly.")
	}
	g.Unit = -1
	if g.convertUnit() != "" {
		t.Error("Grid did not return empty string for improper Unit constant.")
	}
}

func TestConvertPageSize(t *testing.T) {
	g := Grid{
		PageSize: PageSizeLetter,
	}
	if g.convertPageSize() != "Letter" {
		t.Error("Grid did not interpret PageSize constant correctly.")
	}
	g.PageSize = -1
	if g.convertPageSize() != "" {
		t.Error("Grid did not return empty string for improper PageSize constant.")
	}
}

func TestGetCell(t *testing.T) {
	g := newGrid()
	b := Block{5, 2}
	c := g.GetCell(b)
	if c.Width != 218.0 {
		t.Errorf("Grid did not return cell with correct Width. Got %.1f expected %.1f", c.Width, 218.0)
	}
	if c.Height != 24.0 {
		t.Errorf("Grid did not return cell with correct Height. Got %.1f expected %.1f", c.Height, 24.0)
	}
}

func TestGetPoint(t *testing.T) {
	g := newGrid()
	p := g.GetPoint(3, 3)
	if p.X != 128.0 {
		t.Errorf("Grid did not return point with correct X field. Got %.1f expected %.1f", p.X, 128.0)
	}
	if p.Y != 60.0 {
		t.Errorf("Grid did not return point with correct Y field. Got %.1f expected %.1f", p.Y, 60.0)
	}

}
