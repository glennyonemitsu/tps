package tps

import (
	"github.com/jung-kurt/gofpdf"
)

type Report struct {
	Grid             Grid
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
