package parser

import "fmt"

type Registry struct {
	byLang map[string]Parser
}

func NewRegistry(parsers ...Parser) *Registry {
	m := make(map[string]Parser, len(parsers))
	for _, p := range parsers {
		m[p.Language()] = p
	}
	return &Registry{byLang: m}
}

func (r *Registry) Get(lang string) (Parser, error) {
	p, ok := r.byLang[lang]
	if !ok {
		return nil, fmt.Errorf("unsupported language: %s", lang)
	}
	return p, nil
}
