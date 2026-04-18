package parser

import "github.com/melekabbassi/diagramr/internal/graph"

type Parser interface {
	Parse(rootPath string, opts Options) (*graph.DiagramGraph, error)
	Language() string
	Extensions() []string
}

type Options struct {
	IncludePrivate bool
	Exclude        []string // glob patterns to skip
	MaxDepth       int
}
