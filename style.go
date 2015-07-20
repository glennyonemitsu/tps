package tps

// Style is a specification of the content visuals. All content placement
// requires a style name, and cannot be provided dynamically.
type Style struct {
	FontFamily string
	FontStyle  string
	FontSize   float64
	Alignment  int
}

func (s *Style) convertAlignment() string {
	val := ""
	for alignmentConst, stringVal := range alignment {
		if ok := s.Alignment & alignmentConst; ok > 0 {
			val += stringVal
		}
	}
	return val
}
