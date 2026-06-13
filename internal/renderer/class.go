package renderer

import (
	"strings"

	"github.com/melekabbassi/diagramr/internal/graph"
)

type ClassRenderer struct {}

func NewClassRenderer() *ClassRenderer { return &ClassRenderer{} }
func (r *ClassRenderer) Type() string { return "class" }

func (r *ClassRenderer) Render(g *graph.DiagramGraph, opts Options) (string, error) {
	_ = opts
	var b strings.Builder
	b.WriteString("classDiagram\n")
	for _, n := range g.Nodes {
		b.WriteString("  class " + n.Label + "\n")
	}
	return b.String(), nil
}
