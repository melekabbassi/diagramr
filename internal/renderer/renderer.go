package renderer

import "github.com/melekabbassi/diagramr/internal/graph"

type Renderer interface {
	Render(g *graph.DiagramGraph, opts Options) (string, error)
	Type() string
}

type Options struct {
	Direction      string   // TB, LR, BT, RL
	Theme          string
	ShowMethods    bool
	ShowFields     bool
	ShowPrivate    bool
	StripSuffixes  []string
	HideNodes      []string
	OnlyNodes      []string // if set, hide all nodes except for these + neighbors
	MaxNodes       int
	GroupByPackage bool
}
