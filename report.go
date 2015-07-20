package tps

import (
	"encoding/base64"
	"fmt"
	"io"
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

	point := r.Grid.GetPoint(x, y)
	cell := r.Grid.GetCell(block)

	r.Pdf.SetFont(style.FontFamily, style.FontStyle, style.FontSize)
	r.Pdf.SetXY(point.X, point.Y)
	r.Pdf.MultiCell(cell.Width, cell.Height, content, "", style.convertAlignment(), false)

	contentLines := strings.Split(content, "\n")
	for _, line := range contentLines {
		stringWidth := r.Pdf.GetStringWidth(line)
		lineCount += int(math.Ceil(stringWidth / cell.Width))
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
//   r.AddStyle("header", "OpenSans", "", 24, AlignmentCenter | AlignmentTop)
//   r.AddStyle("subheader", "OpenSans", "", 18, AlignmentLeft | AlignmentTop)
func (r *Report) AddStyle(
	name string,
	fontFamily string,
	fontStyle string,
	fontSize float64,
	alignment int,
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
//   r.AddStyle("header", "OpenSans-Bold", "", 64, AlignmentTop | AlignmentLeft)
//
// The following encodings are supported:
//
//	 cp1250
// 	 cp1251
// 	 cp1252
// 	 cp1253
// 	 cp1254
// 	 cp1255
// 	 cp1257
// 	 cp1258
// 	 cp874
// 	 iso-8859-1
// 	 iso-8859-11
// 	 iso-8859-15
// 	 iso-8859-16
// 	 iso-8859-2
// 	 iso-8859-4
// 	 iso-8859-5
// 	 iso-8859-7
// 	 iso-8859-9
// 	 koi8-r
// 	 koi8-u
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
	orientation int,
	pageSize int,
	unit int,
	margin float64,
	columnCount int,
	gutterWidth float64,
	lineHeight float64,
) {
	fontPath := r.FontSourcePath
	r.Grid = Grid{
		Orientation: orientation,
		PageSize:    pageSize,
		Unit:        unit,
		ColumnCount: columnCount,
		GutterWidth: gutterWidth,
		Margin:      margin,
		LineHeight:  lineHeight,
	}

	pdf := gofpdf.New(
		r.Grid.convertOrientation(),
		r.Grid.convertUnit(),
		r.Grid.convertPageSize(),
		fontPath,
	)
	r.Pdf = pdf
	r.Pdf.SetMargins(margin, margin, margin)
	pageWidth, pageHeight := r.Pdf.GetPageSize()
	r.Grid.PageWidth = pageWidth
	r.Grid.PageHeight = pageHeight
	r.Grid.CalculateColumns()
}

// SetFontPath tells the Report where to find fonts specified with AddFont().
func (r *Report) SetFontPath(fontSourcePath string) {
	r.FontSourcePath = fontSourcePath
	r.FontCompiledPath = path.Join(fontSourcePath, "_compiled")
	r.PrepareFontCompiledPath()
	r.Pdf.SetFontLocation(r.FontCompiledPath)
}
