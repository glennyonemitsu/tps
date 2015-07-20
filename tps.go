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

const (
	OrientationPortrait = iota
	OrientationLandscape
)

const (
	PageSizeA3 = iota
	PageSizeA4
	PageSizeA5
	PageSizeLetter
	PageSizeLegal
)

const (
	UnitPt = iota
	UnitMm
	UnitCm
	UnitIn
)

const (
	_ = 1 << iota
	AlignLeft
	AlignCenter
	AlignRight
	AlignTop
	AlignMiddle
	AlignBottom
)

var alignment, orientation, pageSize, unit map[int]string

func init() {
	alignment = map[int]string{
		AlignLeft:   "L",
		AlignCenter: "C",
		AlignRight:  "R",
		AlignTop:    "T",
		AlignMiddle: "M",
		AlignBottom: "B",
	}
	orientation = map[int]string{
		OrientationPortrait:  "Portrait",
		OrientationLandscape: "Landscape",
	}
	pageSize = map[int]string{
		PageSizeA3:     "A3",
		PageSizeA4:     "A4",
		PageSizeA5:     "A5",
		PageSizeLetter: "Letter",
		PageSizeLegal:  "Legal",
	}
	unit = map[int]string{
		UnitPt: "pt",
		UnitMm: "mm",
		UnitCm: "cm",
		UnitIn: "in",
	}
}
