package renderer

import (
	"fmt"
	"sort"
	"strings"

	"github.com/melekabbassi/diagramr/internal/graph"
)

type ClassRenderer struct{}

func NewClassRenderer() *ClassRenderer { return &ClassRenderer{} }
func (r *ClassRenderer) Type() string  { return "class" }

func (r *ClassRenderer) Render(g *graph.DiagramGraph, opts Options) (string, error) {
	var b strings.Builder
	b.WriteString("classDiagram\n")

	ids := make([]string, 0, len(g.Nodes))
	for id := range g.Nodes {
		ids = append(ids, id)
	}
	sort.Strings(ids)

	truncated := false
	if opts.MaxNodes > 0 && len(ids) > opts.MaxNodes {
		ids = ids[:opts.MaxNodes]
		truncated = true
	}

	visible := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		visible[id] = struct{}{}
	}

	for _, id := range ids {
		n := g.Nodes[id]
		label := sanitize(n.Label)
		b.WriteString("  class " + label + " {\n")

		if n.Kind == graph.KindInterface {
			b.WriteString("    <<interface>>\n")
		}

		if opts.ShowFields {
			for _, f := range n.Fields {
				if f.Visibility == graph.VisPrivate && !opts.ShowPrivate {
					continue
				}
				prefix := visPrefix(f.Visibility)
				if f.IsEmbedded {
					b.WriteString(fmt.Sprintf("    %s%s\n", prefix, sanitize(f.Type)))
				} else {
					b.WriteString(fmt.Sprintf("    %s%s %s\n", prefix, sanitize(f.Type), f.Name))
				}
			}
		}

		if opts.ShowMethods {
			for _, m := range n.Methods {
				if m.Visibility == graph.VisPrivate && !opts.ShowPrivate {
					continue
				}
				prefix := visPrefix(m.Visibility)
				params := formatParams(m.Params)
				ret := formatReturns(m.Returns)
				if ret == "" {
					b.WriteString(fmt.Sprintf("    %s%s(%s)\n", prefix, m.Name, params))
				} else {
					b.WriteString(fmt.Sprintf("    %s%s(%s) %s\n", prefix, m.Name, params, ret))
				}
			}
		}

		b.WriteString("  }\n")
	}

	for _, e := range g.Edges {
		if _, ok := visible[e.From]; !ok {
			continue
		}
		if _, ok := visible[e.To]; !ok {
			continue
		}
		from := sanitize(nodeLabel(g, e.From))
		to := sanitize(nodeLabel(g, e.To))
		switch e.Relation {
		case graph.RelImplements:
			b.WriteString(fmt.Sprintf("  %s ..|> %s : implements\n", from, to))
		case graph.RelEmbeds:
			b.WriteString(fmt.Sprintf("  %s --|> %s : embeds\n", from, to))
		case graph.RelUses:
			b.WriteString(fmt.Sprintf("  %s --> %s : uses\n", from, to))
		case graph.RelImports:
			// omitted — package-level imports are too noisy in a class diagram
		}
	}

	if truncated {
		b.WriteString(fmt.Sprintf("  %%%% truncated to %d nodes\n", opts.MaxNodes))
	}

	return b.String(), nil
}

func visPrefix(v graph.Visibility) string {
	if v == graph.VisPublic {
		return "+"
	}
	return "-"
}

func sanitize(s string) string {
	return strings.NewReplacer(
		".", "_", " ", "_", "/", "_",
		"[", "", "]", "", "*", "",
		"{", "", "}", "",
	).Replace(s)
}

func formatParams(params []graph.Param) string {
	parts := make([]string, 0, len(params))
	for _, p := range params {
		if p.Name == "" {
			parts = append(parts, sanitize(p.Type))
		} else {
			parts = append(parts, p.Name+" "+sanitize(p.Type))
		}
	}
	return strings.Join(parts, ", ")
}

func formatReturns(returns []string) string {
	sanitized := make([]string, len(returns))
	for i, r := range returns {
		sanitized[i] = sanitize(r)
	}
	return strings.Join(sanitized, ", ")
}

func nodeLabel(g *graph.DiagramGraph, id string) string {
	if n, ok := g.Nodes[id]; ok {
		return n.Label
	}
	parts := strings.Split(id, ".")
	return parts[len(parts)-1]
}
