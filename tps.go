// tps is a package to hide a lot of the complicated details of creating PDFs
// with fpdf. It relies on a grid specification and content placement is similar
// to a spreadsheet.
//
// For example, a grid specification includes the page margins, # of columns,
// the gutter size (the space between columns), and line height. Using this
// tps calculates a grid coordinate system similar to a speadsheet. Content
// placement is then done with the row, column, and width (determined by the #
// of horizontal cells to take up).
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

// Report is the main struct type that holds all information to generate a PDF.
type Report struct {
	Grid             Grid
	Pdf              *gofpdf.Fpdf
	Styles           map[string]Style
	Blocks           map[string]Block
	FontSourcePath   string
	FontCompiledPath string
}

// Style is a specification of the content visuals. All content placement
// requires a style name, and cannot be provided dynamically.
type Style struct {
	FontFamily string
	FontStyle  string
	FontSize   float64
	Alignment  string
}

// Block is a specification of a line placed in the PDF. The Width and Height
// fields are integers which mean they are the multiples of Grid specs. Width
// indicates # of columns plus the gutters between columns, and Height
// indicates multiples of Grid.LineHeight.
type Block struct {
	Width  int
	Height int
}

func NewReport() *Report {
	report := new(Report)
	report.Styles = make(map[string]Style)
	report.Blocks = make(map[string]Block)
	return report
}

// Place a string based on the x, y coordinates on the grid, using the named
// block and style specifications. Returns the # of lines (different from
// Block.Height) taken up by this call to help dynamically place following
// content.
func (r *Report) Content(
	x int,
	y int,
	blockName string,
	styleName string,
	content string,
) (lineCount int, err error) {
	var block Block
	var style Style
	var ok bool

	lineCount = 0

	if block, ok = r.Blocks[blockName]; ok == false {
		err = fmt.Errorf("Could not find block name in Report: %s", blockName)
		return lineCount, err
	}
	if style, ok = r.Styles[styleName]; ok == false {
		err = fmt.Errorf("Could not find style name in Report: %s", styleName)
		return lineCount, err
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

	contentLines := strings.Split(content, "\n")
	for _, line := range contentLines {
		stringWidth := r.Pdf.GetStringWidth(line)
		lineCount += int(math.Ceil(stringWidth / cellWidth))
	}
	lineCount *= block.Height
	return lineCount, nil
}

// AddPage creates new page in the report. The previous page is now set if it
// exists, and all placement will take place in this new page.
func (r *Report) AddPage() {
	r.Pdf.AddPage()
}

// AddStyle adds a new style to use when placing content in this report.
//
// All specs are set. So even small differences will require different styles.
// An example set of styles can look like the following:
//
//   r.AddStyle("header", "OpenSans", "", 24, "LT")
//   r.AddStyle("subheader", "OpenSans", "", 18, "LT")
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

// AddBlock adds a new block specification to use when placing content in this
// report. The width and height are the number of columns and lines of the block
// respectively. The height is the number of lineHeight per line that the
// content placement takes up.
func (r *Report) AddBlock(name string, width, height int) {
	r.Blocks[name] = Block{
		Width:  width,
		Height: height,
	}
}

// AddFont takes a font filename and compiles it into Report.FontCompiledPath
// with the encoding specified. It strips the filename extension and replaces
// it with .json automatically. The extension-less string becomes the name of
// the font family to use with Report.AddStyle(). For example:
//
//   r.AddFont("OpenSans-Bold.ttf", "cp1252")
//   r.AddStyle("header", "OpenSans-Bold", "", 64, "TF")
func (r *Report) AddFont(filename, encoding string) error {
	var err error
	err = r.PrepareFontCompiledPath()
	if err != nil {
		return err
	}

	ext := path.Ext(filename)
	familyName := filename[:len(filename)-len(ext)]

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

// PrepareFontCompiledPath creates the "_compiled" subdirectory.
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

// CompileEncoding creates the encoding map file in Report.FontCompiledPath so
// the underlying Fpdf object can correctly use it to compile fonts.
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

// CompileFont takes a font file in Report.FontSourcePath and converts it into
// .json format if it doesn't exist in Report.FontCompiledPath.
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

// IsCompiledFile checks if the filename exists in Report.FontCompiledPath. This
// has the suffix "File" instead of "Font" like Report.IsSourcedFont() because
// this method might be used to check encoding map files as well.
func (r *Report) IsCompiledFile(filename string) bool {
	_, err := os.Stat(path.Join(r.FontCompiledPath, filename))
	return !os.IsNotExist(err)
}

// IsSourcedFont checks if the font filename exists in Report.FontSourcePath.
func (r *Report) IsSourcedFont(filename string) bool {
	_, err := os.Stat(path.Join(r.FontSourcePath, filename))
	return !os.IsNotExist(err)
}

// SetGrid sets all page and grid related specifications required to place
// content. This must be set before any Content() calls are made.
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

// SetFontPath tells the Report where to find fonts specified with AddFont().
func (r *Report) SetFontPath(fontSourcePath string) {
	r.FontSourcePath = fontSourcePath
	r.FontCompiledPath = path.Join(fontSourcePath, "_compiled")
	r.PrepareFontCompiledPath()
	r.Pdf.SetFontLocation(r.FontCompiledPath)
}
