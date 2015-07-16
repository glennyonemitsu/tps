package tps

import (
	"github.com/jung-kurt/gofpdf"
)

type Report struct {
	ColumnCount      float64
	ColumnWidth      float64
	GutterCount      float64
	GutterWidth      float64
	LineHeight       float64
	Margin           float64
	Pdf              *gofpdf.Fpdf
	Styles           map[string]Style
	Blocks           map[string]Block
	FontSourcePath   string
	FontCompiledPath string
}

type Style struct {
	FontFamily string
	FontStyle  string
	FontSize   float64
	Alignment  string
}

type Block struct {
	Width  int
	Height int
}
