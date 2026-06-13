package studio

import (
	"fmt"

	base "github.com/melekabbassi/diagramr/internal/parser"
	golangparser "github.com/melekabbassi/diagramr/internal/parser/golang"
	"github.com/melekabbassi/diagramr/internal/renderer"
)

// RenderConfig holds everything needed to (re-)parse and render a diagram.
type RenderConfig struct {
	RootPath   string
	Lang       string
	ParseOpts  base.Options
	RenderOpts renderer.Options
}

func renderMermaid(cfg RenderConfig) (string, error) {
	switch cfg.Lang {
	case "go":
		g, err := golangparser.New().Parse(cfg.RootPath, cfg.ParseOpts)
		if err != nil {
			return "", fmt.Errorf("parse: %w", err)
		}
		out, err := renderer.NewClassRenderer().Render(g, cfg.RenderOpts)
		if err != nil {
			return "", fmt.Errorf("render: %w", err)
		}
		return out, nil
	default:
		return "", fmt.Errorf("unsupported language %q", cfg.Lang)
	}
}
