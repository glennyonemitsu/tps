package tps

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"path"
	"strings"

	"github.com/jung-kurt/gofpdf"
)

func NewReport() *Report {
	report := new(Report)
	report.Styles = make(map[string]Style)
	report.Blocks = make(map[string]Block)
	return report
}

func (r *Report) Content(
	x int,
	y int,
	blockName string,
	styleName string,
	content string,
) int {
	var block Block
	var style Style
	var ok bool

	if block, ok = r.Blocks[blockName]; ok == false {
		log.Fatalf("Could not find block name in Report: %s", blockName)
		os.Exit(1)
	}
	if style, ok = r.Styles[styleName]; ok == false {
		log.Fatalf("Could not find style name in Report: %s", styleName)
		os.Exit(2)
	}

	pageX := r.Grid.Margin
	pageX += r.Grid.ColumnWidth * float64(x-1)
	pageX += r.Grid.GutterWidth * float64(x-1)

	pageY := r.Grid.Margin
	pageY += r.Grid.LineHeight * float64(y-1)

	cellWidth := r.Grid.ColumnWidth * float64(block.Width)
	cellWidth += r.Grid.GutterWidth * float64(block.Width-1)
	cellHeight := r.Grid.LineHeight * float64(block.Height)

	r.Pdf.SetFont(style.FontFamily, style.FontStyle, style.FontSize)
	r.Pdf.SetXY(pageX, pageY)
	r.Pdf.MultiCell(cellWidth, cellHeight, content, "", style.Alignment, false)

	lineCount := 0
	contentLines := strings.Split(content, "\n")
	for _, line := range contentLines {
		stringWidth := r.Pdf.GetStringWidth(line)
		lineCount += int(math.Ceil(stringWidth / cellWidth))
	}
	lineCount *= block.Height
	return lineCount
}

func (r *Report) AddPage() {
	r.Pdf.AddPage()
}

func (r *Report) AddStyle(
	name string,
	fontFamily string,
	fontStyle string,
	fontSize float64,
	alignment string,
) {
	r.Styles[name] = Style{
		FontFamily: fontFamily,
		FontStyle:  fontStyle,
		FontSize:   fontSize,
		Alignment:  alignment,
	}
}

func (r *Report) AddBlock(name string, width, height int) {
	r.Blocks[name] = Block{
		Width:  width,
		Height: height,
	}
}

func (r *Report) AddFont(filename, encoding string) error {
	var err error
	err = r.PrepareFontCompiledPath()
	if err != nil {
		return err
	}

	familyName := stripExt(filename)

	// auto compiles
	if path.Ext(filename) == ".json" {
		if r.IsCompiledFile(filename) {
			r.Pdf.AddFont(familyName, "", filename)
		} else {
			return fmt.Errorf("Cache font file not found: %s", filename)
		}
	} else {
		if r.IsSourcedFont(filename) {
			compiledFilename, err := r.CompileFont(filename, encoding)
			if err != nil {
				return fmt.Errorf("Could not compile font: %v", err)
			}
			r.Pdf.AddFont(familyName, "", compiledFilename)
		} else {
			return fmt.Errorf("Source font file not found: %s", filename)
		}
	}
	return nil
}

func (r *Report) PrepareFontCompiledPath() error {
	if _, err := os.Stat(path.Join(r.FontCompiledPath)); os.IsNotExist(err) {
		err = os.MkdirAll(r.FontCompiledPath, os.ModeDir)
		if err != nil {
			return err
		}
		return os.Chmod(r.FontCompiledPath, 0775)
	}
	return nil
}

func (r *Report) CompileEncoding(encoding string) (filename string, err error) {
	filename = path.Join(r.FontCompiledPath, encoding+".map")
	if r.IsCompiledFile(filename) {
		return
	}
	if data, ok := encodings[encoding]; ok {
		file, err := os.Create(filename)
		if err != nil {
			err = fmt.Errorf("Could not open file to compile encoding file: %v", err)
			return filename, err
		}
		defer file.Close()
		reader := strings.NewReader(data)
		decoder := base64.NewDecoder(base64.StdEncoding, reader)
		_, err = io.Copy(file, decoder)
		if err != nil {
			err = fmt.Errorf("Encoding failed to copy to file: %v", err)
			return filename, err
		}
	} else {
		err = fmt.Errorf("Encoding not supported: %s", encoding)
	}
	return
}

func (r *Report) CompileFont(filename, encoding string) (string, error) {
	fontFilename := path.Join(r.FontSourcePath, filename)
	// replacing ext with json
	extLen := len(path.Ext(filename))
	compiledFilename := filename[:len(filename)-extLen] + ".json"
	encodingFilename, err := r.CompileEncoding(encoding)
	if err != nil {
		return compiledFilename, err
	}
	err = gofpdf.MakeFont(fontFilename, encodingFilename, r.FontCompiledPath, nil, true)
	return compiledFilename, err
}

func (r *Report) IsCompiledFile(filename string) bool {
	_, err := os.Stat(path.Join(r.FontCompiledPath, filename))
	return !os.IsNotExist(err)
}

func (r *Report) IsSourcedFont(filename string) bool {
	_, err := os.Stat(path.Join(r.FontSourcePath, filename))
	return !os.IsNotExist(err)
}

func (r *Report) SetGrid(
	orientation string,
	size string,
	units string,
	margin float64,
	columnCount int,
	gutterWidth float64,
	lineHeight float64,
) {
	fontPath := r.FontSourcePath
	pdf := gofpdf.New(orientation, units, size, fontPath)
	r.Pdf = pdf
	r.Pdf.SetMargins(margin, margin, margin)
	pageWidth, pageHeight := r.Pdf.GetPageSize()
	r.Grid = Grid{
		Orientation: orientation,
		Size:        size,
		Units:       units,
		ColumnCount: columnCount,
		GutterWidth: gutterWidth,
		PageWidth:   pageWidth,
		PageHeight:  pageHeight,
		Margin:      margin,
		LineHeight:  lineHeight,
	}
	r.Grid.CalculateColumns()
}

func (r *Report) SetFontPath(fontSourcePath string) {
	r.FontSourcePath = fontSourcePath
	r.FontCompiledPath = path.Join(fontSourcePath, "_compiled")
	r.PrepareFontCompiledPath()
	r.Pdf.SetFontLocation(r.FontCompiledPath)
}
