package tps

import (
	"testing"
)

func TestAddBlock(t *testing.T) {
	r := NewReport()
	r.AddBlock("test", 1, 2)
	e := Block{1, 2}
	if b := r.Blocks["test"]; b != e {
		t.Errorf("AddBlock did not store block correctly. Got %v expected %v", b, e)
	}

	r.AddBlock("test", 3, 4)
	e = Block{3, 4}
	if b := r.Blocks["test"]; b != e {
		t.Errorf("AddBlock did not overwrite block correctly. Got %v expected %v", b, e)
	}

	r.AddBlock("new test", 5, 6)
	e = Block{5, 6}
	if b := r.Blocks["new test"]; b != e {
		t.Errorf("AddBlock did not store block correctly. Got %v expected %v", b, e)
	}
}

func TestAddStyle(t *testing.T) {
	r := NewReport()
	a := AlignLeft | AlignTop
	r.AddStyle("test", "foo", "", 12, a)
	e := Style{"foo", "", 12, a}
	if s := r.Styles["test"]; s != e {
		t.Errorf("AddStyle did not store style correctly. Got %v expected %v", s, e)
	}

	r.AddStyle("test", "foo bar", "", 24, a)
	e = Style{"foo bar", "", 24, a}
	if s := r.Styles["test"]; s != e {
		t.Errorf("AddStyle did not overwrite style correctly. Got %v expected %v", s, e)
	}

	r.AddStyle("new test", "foo bar", "", 24, a)
	e = Style{"foo bar", "", 24, a}
	if s := r.Styles["new test"]; s != e {
		t.Errorf("AddStyle did not store style correctly. Got %v expected %v", s, e)
	}
}

func TestSetGrid(t *testing.T) {
	r := NewReport()
	r.SetGrid(OrientationPortrait, PageSizeLetter, UnitPt, 36.0, 12, 12.0, 12.0)
	g := r.Grid
	e := Grid{
		ColumnCount: 12,
		ColumnWidth: 34.0,
		GutterCount: 11,
		GutterWidth: 12.0,
		LineHeight:  12.0,
		Margin:      36.0,
		Orientation: OrientationPortrait,
		PageWidth:   612.0,
		PageHeight:  792.0,
		PageSize:    PageSizeLetter,
		Unit:        UnitPt,
	}
	if g != e {
		t.Errorf("SetGrid did not set correct Grid. Got %v expected %v", g, e)
	}
	if r.Pdf == nil {
		t.Errorf("SetGrid did not initialize Pdf")
	}
}
