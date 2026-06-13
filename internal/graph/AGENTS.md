<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-06-13 | Updated: 2026-06-13 -->

# internal/graph

## Purpose
The intermediate representation (IR) that all parsers produce and all renderers consume. Defines the `DiagramGraph` and its constituent types — nodes, edges, fields, methods, and metadata.

## Key Files

| File | Description |
|------|-------------|
| `model.go` | All graph types: `DiagramGraph`, `Node`, `Edge`, `Method`, `Field`, `Param`, `Metadata`, and their enum constants |
| `model_test.go` | Snapshot-style tests verifying the JSON serialization shape is stable |

## For AI Agents

### Working In This Directory
- `DiagramGraph.Nodes` is a `map[string]*Node` keyed by fully-qualified ID (`pkgName.TypeName`).
- Node kinds: `struct`, `interface`, `service` (inferred by naming convention).
- Edge relations: `embeds`, `implements`, `uses`, `imports`.
- Visibility is `public` (exported, uppercase) or `private` (unexported, lowercase).
- **Do not** add parser or renderer logic here — this package must remain a pure data model with no dependencies on other internal packages.

### Testing Requirements
- Add serialization tests to `model_test.go` when adding new fields, to guard against accidental JSON shape changes.
- `go test ./internal/graph/... -race`

## Dependencies

### Internal
- None — this package is a leaf with no internal imports.

<!-- MANUAL: Any manually added notes below this line are preserved on regeneration -->
