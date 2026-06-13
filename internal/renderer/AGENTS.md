<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-06-13 | Updated: 2026-06-13 -->

# internal/renderer

## Purpose
Defines the `Renderer` interface and rendering options. Converts a `graph.DiagramGraph` into a format-specific string (e.g., Mermaid class diagram syntax). The class diagram renderer lives here.

## Key Files

| File | Description |
|------|-------------|
| `renderer.go` | `Renderer` interface (`Render`, `Type`) and `Options` struct with layout/filter settings |
| `class.go` | Mermaid class diagram renderer — converts graph nodes and edges to Mermaid `classDiagram` syntax |

## For AI Agents

### Working In This Directory
- `Options` controls: flow direction (`TB`/`LR`/`BT`/`RL`), theme, visibility filters (`ShowMethods`, `ShowFields`, `ShowPrivate`), node filtering (`HideNodes`, `OnlyNodes`), node limit (`MaxNodes`), and package grouping.
- To add a new renderer (e.g., PlantUML), implement the `Renderer` interface in a new file.
- Renderers produce strings; writing to disk is handled by `internal/output`.

### Testing Requirements
- Test renderers with known `DiagramGraph` inputs and assert the output string matches expected diagram syntax.
- `go test ./internal/renderer/... -race`

## Dependencies

### Internal
- `internal/graph` — `DiagramGraph` is the sole input

<!-- MANUAL: Any manually added notes below this line are preserved on regeneration -->
