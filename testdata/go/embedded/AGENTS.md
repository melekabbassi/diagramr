<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-06-13 | Updated: 2026-06-13 -->

# testdata/go/embedded

## Purpose
Fixture for Go struct embedding. Tests that the parser emits `RelEmbeds` edges between structs that use anonymous (embedded) fields.

## Key Files

| File | Description |
|------|-------------|
| `embedded.go` | Package `embedded` — structs with embedded types to exercise `emitEmbeddingEdges` |

## For AI Agents

### Working In This Directory
- Embedded fields have no `Names` in the AST `FieldList` — the parser identifies them via `len(f.Names) == 0`.
- Keep this fixture focused on embedding only; do not mix in interface satisfaction or other relationships.

<!-- MANUAL: Any manually added notes below this line are preserved on regeneration -->
