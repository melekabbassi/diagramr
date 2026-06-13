package renderer_test

import (
	"strings"
	"testing"

	"github.com/melekabbassi/diagramr/internal/graph"
	"github.com/melekabbassi/diagramr/internal/renderer"
)

func TestClassRenderer_Render(t *testing.T) {
	tests := []struct {
		name    string
		g       *graph.DiagramGraph
		opts    renderer.Options
		want    []string // substrings that must appear in output
		wantErr bool
	}{
		{
			name: "empty graph",
			g:    &graph.DiagramGraph{Nodes: map[string]*graph.Node{}, Edges: []graph.Edge{}},
			opts: renderer.Options{},
			want: []string{"classDiagram"},
		},
		{
			name: "single struct with public field and method",
			g: &graph.DiagramGraph{
				Nodes: map[string]*graph.Node{
					"MyStruct": {
						Kind: graph.KindStruct,
						Fields: []graph.Field{
							{Name: "ID", Type: "int", Visibility: graph.VisPublic},
						},
						Methods: []graph.Method{
							{Name: "Run", Visibility: graph.VisPublic},
						},
					},
				},
				Edges: []graph.Edge{},
			},
			opts: renderer.Options{ShowFields: true, ShowMethods: true},
			want: []string{"class MyStruct", "+int ID", "+Run()"},
		},
		{
			name: "interface node",
			g: &graph.DiagramGraph{
				Nodes: map[string]*graph.Node{
					"MyIface": {Kind: graph.KindInterface},
				},
				Edges: []graph.Edge{},
			},
			opts: renderer.Options{ShowMethods: true},
			want: []string{"<<interface>>"},
		},
		{
			name: "direction is emitted",
			g:    &graph.DiagramGraph{Nodes: map[string]*graph.Node{}, Edges: []graph.Edge{}},
			opts: renderer.Options{Direction: "LR"},
			want: []string{"direction LR"},
		},
		{
			name: "private fields hidden without ShowPrivate",
			g: &graph.DiagramGraph{
				Nodes: map[string]*graph.Node{
					"S": {
						Kind: graph.KindStruct,
						Fields: []graph.Field{
							{Name: "secret", Type: "string", Visibility: graph.VisPrivate},
						},
					},
				},
				Edges: []graph.Edge{},
			},
			opts: renderer.Options{ShowFields: true, ShowPrivate: false},
			want: []string{"class S"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := renderer.NewClassRenderer()
			got, err := r.Render(tt.g, tt.opts)
			if err != nil {
				if !tt.wantErr {
					t.Fatalf("Render() unexpected error: %v", err)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("Render() succeeded unexpectedly")
			}
			for _, sub := range tt.want {
				if !strings.Contains(got, sub) {
					t.Errorf("Render() output missing %q\ngot:\n%s", sub, got)
				}
			}
		})
	}
}
