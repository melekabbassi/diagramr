<!-- Parent: ../AGENTS.md -->
<!-- Generated: 2026-06-13 | Updated: 2026-06-13 -->

# .github/workflows

## Purpose
GitHub Actions CI workflow definitions. Currently a single `ci.yml` that builds, tests, and vets the project on all three major OSes.

## Key Files

| File | Description |
|------|-------------|
| `ci.yml` | CI pipeline — matrix over ubuntu/macos/windows; runs `go build`, `go test -race`, `go vet` |

## For AI Agents

### Working In This Directory
- The CI matrix uses Go 1.22 (`actions/setup-go@v5`) — match this version locally when debugging CI failures.
- All three steps (`build`, `test`, `vet`) must pass on all three OSes before merging.
- Add new jobs as separate entries under `jobs:` in `ci.yml`; keep the build job unchanged.
- YAML is indentation-sensitive — validate with `yamllint` or the GitHub Actions VSCode extension before pushing.

### Testing Requirements
- Use `act` to run workflows locally: `act push -j ci`

<!-- MANUAL: Any manually added notes below this line are preserved on regeneration -->
