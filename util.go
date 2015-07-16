package tps

import (
	"path"
)

func stripExt(filename string) string {
	ext := path.Ext(filename)
	return filename[:len(filename)-len(ext)]
}
