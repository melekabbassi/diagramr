<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-06-13 | Updated: 2026-06-13 -->

# internal/output

## Purpose
Output writers that serialize a rendered diagram (string or binary) to files or stdout. Each file handles one format. Currently stubs — the real serialization logic is to be implemented.

## Key Files

| File | Description |
|------|-------------|
| `writer.go` | `Write(format, content string) error` — dispatcher stub routing to format-specific writers |
| `json.go` | `RenderJSON(v any) (string, error)` — JSON serialization of the graph or rendered output |
| `mermaid.go` | Mermaid text output writer (stub) |
| `svg.go` | SVG output writer (stub) |
| `png.go` | PNG output writer (stub) |

## For AI Agents

### Working In This Directory
- Implement format writers by filling in the stub functions; `writer.go` should dispatch based on the `format` string.
- JSON output should marshal `graph.DiagramGraph` using `encoding/json`.
- SVG and PNG likely require rendering the Mermaid text first via an external tool (e.g., mmdc) or a Go Mermaid library.

### Testing Requirements
- Unit test each format writer with a known `DiagramGraph` input and assert output correctness.
- `go test ./internal/output/... -race`

## Dependencies

### Internal
- `internal/graph` — `DiagramGraph` is the primary input type

<!-- MANUAL: Any manually added notes below this line are preserved on regeneration -->
