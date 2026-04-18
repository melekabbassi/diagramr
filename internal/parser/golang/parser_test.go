package golang

import (
	"path/filepath"
	"testing"

	"github.com/melekabbassi/diagramr/internal/graph"
	base "github.com/melekabbassi/diagramr/internal/parser"
)

func TestParseSimpleFixture(t *testing.T) {
	p := New()
	root := fixturePath(t, "simple")

	g, err := p.Parse(root, base.Options{IncludePrivate: false})
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	user, ok := g.Nodes["simple.User"]
	if !ok {
		t.Fatalf("expected simple.User node to exist")
	}

	if user.Kind != graph.KindStruct {
		t.Fatalf("expected struct kind, got %q", user.Kind)
	}

	if len(user.Methods) != 1 {
		t.Fatalf("expected 1 public method, got %d", len(user.Methods))
	}

	if user.Methods[0].Name != "FullName" || !user.Methods[0].IsPointer {
		t.Fatalf("unexpected method extraction: %+v", user.Methods[0])
	}

	if containsField(user.Fields, "secret") {
		t.Fatalf("private field should be excluded when IncludePrivate=false")
	}

	if !containsField(user.Fields, "Address") {
		t.Fatalf("expected Address field in extracted fields")
	}
}

func TestParseEmbeddedFixture(t *testing.T) {
	p := New()
	root := fixturePath(t, "embedded")

	g, err := p.Parse(root, base.Options{IncludePrivate: true})
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	if !containsEdge(g.Edges, "embedded.User", "embedded.Base", graph.RelEmbeds) {
		t.Fatalf("expected embeds edge embedded.User -> embedded.Base")
	}
}

func TestParseInterfacesFixture(t *testing.T) {
	p := New()
	root := fixturePath(t, "interfaces")

	g, err := p.Parse(root, base.Options{IncludePrivate: true})
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	if !containsEdge(g.Edges, "interfaces.FileReader", "interfaces.Reader", graph.RelImplements) {
		t.Fatalf("expected implements edge interfaces.FileReader -> interfaces.Reader")
	}
}

func TestParseMultifileFixture(t *testing.T) {
	p := New()
	root := fixturePath(t, "multifile")

	g, err := p.Parse(root, base.Options{IncludePrivate: true})
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	user, ok := g.Nodes["multifile.User"]
	if !ok {
		t.Fatalf("expected multifile.User node to exist")
	}

	if len(user.Methods) != 1 || user.Methods[0].Name != "Rename" {
		t.Fatalf("expected Rename method from separate file, got %+v", user.Methods)
	}
}

func fixturePath(t *testing.T, name string) string {
	t.Helper()
	return filepath.Join("..", "..", "..", "testdata", "go", name)
}

func containsField(fields []graph.Field, name string) bool {
	for _, f := range fields {
		if f.Name == name {
			return true
		}
	}
	return false
}

func containsEdge(edges []graph.Edge, from, to string, relation graph.Relation) bool {
	for _, e := range edges {
		if e.From == from && e.To == to && e.Relation == relation {
			return true
		}
	}
	return false
}
