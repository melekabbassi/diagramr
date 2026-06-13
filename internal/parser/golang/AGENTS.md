<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-06-13 | Updated: 2026-06-13 -->

# internal/parser/golang

## Purpose
Go AST-based parser that walks a source tree, extracts struct and interface types, and infers relationships (embedding, implementation, field usage, imports). The primary implemented parser in diagramr.

## Key Files

| File | Description |
|------|-------------|
| `parser.go` | Full parser implementation — `Parser` struct, `Parse()`, and all AST helper functions |
| `parser_test.go` | Fixture-driven tests using `testdata/go/` scenarios |

## For AI Agents

### Working In This Directory
- `Parse()` pipeline: `sourceDirs` (walk for non-test `.go` files) → `parser.ParseDir` per dir → `collectTypes` → `collectMethodsAndImports` → emit edge functions → sort edges.
- Node ID format is `pkgName.TypeName` (e.g., `simple.User`).
- Three edge-emit passes run per package: `emitEmbeddingEdges`, `emitUsageEdges`, `emitImplementsEdges`.
- `emitImplementsEdges` is structural (duck-typed): a struct implements an interface if it has all the interface's methods by name — no explicit `implements` declaration needed.
- Test files (`_test.go`) and non-`.go` files are intentionally skipped by `sourceDirs`.
- `opts.IncludePrivate=false` filters unexported fields and methods from both nodes and edges.

### Testing Requirements
- Tests point at subdirectories of `testdata/go/` via relative paths.
- Add a new fixture dir and a corresponding `Test*` function for each new parsing scenario.
- `go test ./internal/parser/golang/... -race`

### Common Patterns
- `typeString(ast.Expr) string` — recursive type stringification for fields and params.
- `baseTypeName(t string)` — strips `*`, `[]`, `map[K]` prefixes to get the raw type name for edge lookup.
- Edges are deduplicated via `edgeSet` (`relation|from|to` key) before appending to the graph.

## Dependencies

### Internal
- `internal/graph` — produces `*graph.DiagramGraph`
- `internal/parser` (base) — implements `parser.Parser` interface; uses `parser.Options`

### External
- `go/ast`, `go/parser`, `go/token` — standard library Go AST packages

<!-- MANUAL: Any manually added notes below this line are preserved on regeneration -->
