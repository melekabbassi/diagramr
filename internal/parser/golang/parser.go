package golang

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"unicode"

	"github.com/melekabbassi/diagramr/internal/graph"
	base "github.com/melekabbassi/diagramr/internal/parser"
)

type Parser struct{}

func New() *Parser { return &Parser{} }

func (p *Parser) Language() string { return "go" }

func (p *Parser) Extensions() []string { return []string{".go"} }

func (p *Parser) Parse(rootPath string, opts base.Options) (*graph.DiagramGraph, error) {
	g := &graph.DiagramGraph{
		Nodes: make(map[string]*graph.Node),
		Edges: make([]graph.Edge, 0),
		Metadata: graph.Metadata{
			Language: "go",
			RootPath: rootPath,
			ParsedAt: time.Now().UTC().Format(time.RFC3339),
		},
	}

	edgeSet := make(map[string]struct{})

	dirs, err := sourceDirs(rootPath)
	if err != nil {
		return nil, fmt.Errorf("scan source dirs: %w", err)
	}

	for _, dir := range dirs {
		fset := token.NewFileSet()
		pkgs, err := parser.ParseDir(fset, dir, nil, 0)
		if err != nil {
			return nil, fmt.Errorf("parse dir %q: %w", dir, err)
		}

		for pkgName, pkg := range pkgs {
			for filePath, f := range pkg.Files {
				collectTypes(fset, filePath, pkgName, f, g, opts)
			}

			for _, f := range pkg.Files {
				collectMethodsAndImports(pkgName, f, g, opts, edgeSet)
			}

			emitEmbeddingEdges(pkgName, g, edgeSet)
			emitUsageEdges(pkgName, g, edgeSet)
			emitImplementsEdges(pkgName, g, edgeSet)
		}
	}

	sort.Slice(g.Edges, func(i, j int) bool {
		if g.Edges[i].From == g.Edges[j].From {
			if g.Edges[i].To == g.Edges[j].To {
				return g.Edges[i].Relation < g.Edges[j].Relation
			}
			return g.Edges[i].To < g.Edges[j].To
		}
		return g.Edges[i].From < g.Edges[j].From
	})

	return g, nil
}

func sourceDirs(rootPath string) ([]string, error) {
	dirSet := make(map[string]struct{})
	err := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, "_test.go") || filepath.Ext(path) != ".go" {
			return nil
		}
		dirSet[filepath.Dir(path)] = struct{}{}
		return nil
	})
	if err != nil {
		return nil, err
	}

	dirs := make([]string, 0, len(dirSet))
	for dir := range dirSet {
		dirs = append(dirs, dir)
	}
	sort.Strings(dirs)
	return dirs, nil
}

func collectTypes(fset *token.FileSet, filePath, pkgName string, f *ast.File, g *graph.DiagramGraph, opts base.Options) {
	ast.Inspect(f, func(n ast.Node) bool {
		node, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		var kind graph.NodeKind
		switch node.Type.(type) {
		case *ast.StructType:
			kind = graph.KindStruct
		case *ast.InterfaceType:
			kind = graph.KindInterface
		default:
			return true
		}

		id := buildNodeID(pkgName, node.Name.Name)
		g.Nodes[id] = &graph.Node{
			ID:         id,
			Label:      node.Name.Name,
			Kind:       kind,
			Methods:    make([]graph.Method, 0),
			Fields:     make([]graph.Field, 0),
			Package:    pkgName,
			SourceFile: filePath,
			SourceLine: fset.Position(node.Pos()).Line,
		}

		switch t := node.Type.(type) {
		case *ast.StructType:
			addStructFields(t, g.Nodes[id], opts)
		case *ast.InterfaceType:
			addInterfaceMethods(t, g.Nodes[id], opts)
		}

		return true
	})
}

func collectMethodsAndImports(pkgName string, f *ast.File, g *graph.DiagramGraph, opts base.Options, edgeSet map[string]struct{}) {
	for _, imp := range f.Imports {
		importPath := strings.Trim(imp.Path.Value, "\"")
		addEdge(g, edgeSet, graph.Edge{From: pkgName, To: importPath, Relation: graph.RelImports})
	}

	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Recv == nil {
			continue
		}

		recvName, isPointer := receiverType(fn.Recv)
		if recvName == "" {
			continue
		}

		nodeID := buildNodeID(pkgName, recvName)
		n, ok := g.Nodes[nodeID]
		if !ok {
			continue
		}

		m := graph.Method{
			Name:       fn.Name.Name,
			Params:     paramsFromFieldList(fn.Type.Params),
			Returns:    returnsFromFieldList(fn.Type.Results),
			Visibility: visibilityForName(fn.Name.Name),
			IsPointer:  isPointer,
		}

		if !opts.IncludePrivate && m.Visibility == graph.VisPrivate {
			continue
		}

		n.Methods = append(n.Methods, m)
	}
}

func addStructFields(st *ast.StructType, n *graph.Node, opts base.Options) {
	if st.Fields == nil {
		return
	}

	for _, f := range st.Fields.List {
		if len(f.Names) == 0 {
			n.Fields = append(n.Fields, graph.Field{
				Name:       "",
				Type:       typeString(f.Type),
				Visibility: graph.VisPublic,
				IsEmbedded: true,
				Tag:        fieldTag(f.Tag),
			})
			continue
		}

		for _, name := range f.Names {
			field := graph.Field{
				Name:       name.Name,
				Type:       typeString(f.Type),
				Visibility: visibilityForName(name.Name),
				IsEmbedded: false,
				Tag:        fieldTag(f.Tag),
			}

			if !opts.IncludePrivate && field.Visibility == graph.VisPrivate {
				continue
			}

			n.Fields = append(n.Fields, field)
		}
	}
}

func addInterfaceMethods(it *ast.InterfaceType, n *graph.Node, opts base.Options) {
	if it.Methods == nil {
		return
	}

	for _, m := range it.Methods.List {
		if len(m.Names) == 0 {
			continue
		}

		fnType, ok := m.Type.(*ast.FuncType)
		if !ok {
			continue
		}

		for _, name := range m.Names {
			method := graph.Method{
				Name:       name.Name,
				Params:     paramsFromFieldList(fnType.Params),
				Returns:    returnsFromFieldList(fnType.Results),
				Visibility: visibilityForName(name.Name),
				IsPointer:  false,
			}

			if !opts.IncludePrivate && method.Visibility == graph.VisPrivate {
				continue
			}

			n.Methods = append(n.Methods, method)
		}
	}
}

func emitEmbeddingEdges(pkgName string, g *graph.DiagramGraph, edgeSet map[string]struct{}) {
	for _, n := range g.Nodes {
		if n.Package != pkgName || n.Kind != graph.KindStruct {
			continue
		}

		for _, field := range n.Fields {
			if !field.IsEmbedded {
				continue
			}

			typeName := baseTypeName(field.Type)
			targetID := buildNodeID(pkgName, typeName)
			if _, ok := g.Nodes[targetID]; !ok {
				continue
			}

			addEdge(g, edgeSet, graph.Edge{From: n.ID, To: targetID, Relation: graph.RelEmbeds})
		}
	}
}

func emitUsageEdges(pkgName string, g *graph.DiagramGraph, edgeSet map[string]struct{}) {
	for _, n := range g.Nodes {
		if n.Package != pkgName || n.Kind != graph.KindStruct {
			continue
		}

		for _, field := range n.Fields {
			if field.IsEmbedded {
				continue
			}

			typeName := baseTypeName(field.Type)
			targetID := buildNodeID(pkgName, typeName)
			if _, ok := g.Nodes[targetID]; !ok {
				continue
			}

			addEdge(g, edgeSet, graph.Edge{From: n.ID, To: targetID, Relation: graph.RelUses})
		}
	}
}

func emitImplementsEdges(pkgName string, g *graph.DiagramGraph, edgeSet map[string]struct{}) {
	interfaces := make([]*graph.Node, 0)
	structs := make([]*graph.Node, 0)

	for _, n := range g.Nodes {
		if n.Package != pkgName {
			continue
		}
		if n.Kind == graph.KindInterface {
			interfaces = append(interfaces, n)
		}
		if n.Kind == graph.KindStruct || n.Kind == graph.KindService {
			structs = append(structs, n)
		}
	}

	for _, s := range structs {
		sMethods := make(map[string]struct{}, len(s.Methods))
		for _, m := range s.Methods {
			sMethods[m.Name] = struct{}{}
		}

		for _, i := range interfaces {
			if len(i.Methods) == 0 {
				continue
			}

			all := true
			for _, im := range i.Methods {
				if _, ok := sMethods[im.Name]; !ok {
					all = false
					break
				}
			}
			if all {
				addEdge(g, edgeSet, graph.Edge{From: s.ID, To: i.ID, Relation: graph.RelImplements})
			}
		}
	}
}

func addEdge(g *graph.DiagramGraph, edgeSet map[string]struct{}, e graph.Edge) {
	key := string(e.Relation) + "|" + e.From + "|" + e.To
	if _, ok := edgeSet[key]; ok {
		return
	}
	edgeSet[key] = struct{}{}
	g.Edges = append(g.Edges, e)
}

func buildNodeID(pkgName, typeName string) string {
	return pkgName + "." + typeName
}

func visibilityForName(name string) graph.Visibility {
	if name == "" {
		return graph.VisPrivate
	}
	r := []rune(name)[0]
	if unicode.IsUpper(r) {
		return graph.VisPublic
	}
	return graph.VisPrivate
}

func receiverType(recv *ast.FieldList) (string, bool) {
	if recv == nil || len(recv.List) == 0 {
		return "", false
	}

	switch t := recv.List[0].Type.(type) {
	case *ast.StarExpr:
		if ident, ok := t.X.(*ast.Ident); ok {
			return ident.Name, true
		}
	case *ast.Ident:
		return t.Name, false
	}

	return "", false
}

func paramsFromFieldList(fl *ast.FieldList) []graph.Param {
	if fl == nil {
		return nil
	}

	params := make([]graph.Param, 0)
	for _, field := range fl.List {
		t := typeString(field.Type)
		if len(field.Names) == 0 {
			params = append(params, graph.Param{Name: "", Type: t})
			continue
		}
		for _, name := range field.Names {
			params = append(params, graph.Param{Name: name.Name, Type: t})
		}
	}
	return params
}

func returnsFromFieldList(fl *ast.FieldList) []string {
	if fl == nil {
		return nil
	}

	returns := make([]string, 0, len(fl.List))
	for _, field := range fl.List {
		returns = append(returns, typeString(field.Type))
	}
	return returns
}

func typeString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + typeString(t.X)
	case *ast.ArrayType:
		return "[]" + typeString(t.Elt)
	case *ast.MapType:
		return "map[" + typeString(t.Key) + "]" + typeString(t.Value)
	case *ast.SelectorExpr:
		return typeString(t.X) + "." + t.Sel.Name
	case *ast.InterfaceType:
		return "interface{}"
	default:
		return "unknown"
	}
}

func baseTypeName(t string) string {
	clean := strings.TrimPrefix(t, "*")
	clean = strings.TrimPrefix(clean, "[]")
	if strings.HasPrefix(clean, "map[") {
		idx := strings.Index(clean, "]")
		if idx != -1 && idx+1 < len(clean) {
			clean = clean[idx+1:]
		}
	}
	if strings.Contains(clean, ".") {
		parts := strings.Split(clean, ".")
		clean = parts[len(parts)-1]
	}
	return clean
}

func fieldTag(tag *ast.BasicLit) string {
	if tag == nil {
		return ""
	}
	return strings.Trim(tag.Value, "`")
}
