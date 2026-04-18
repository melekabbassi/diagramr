package parser

import "slices"

import "path/filepath"

func MatchExtension(path string, exts []string) bool {
	ext := filepath.Ext(path)
	return slices.Contains(exts, ext)
}
