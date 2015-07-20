package tps

// Block is a specification of a line placed in the PDF. The Width and Height
// fields are integers which mean they are the multiples of Grid specs. Width
// indicates # of columns plus the gutters between columns, and Height
// indicates multiples of Grid.LineHeight.
type Block struct {
	Width  int
	Height int
}
