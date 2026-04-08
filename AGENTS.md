# AGENTS.md — Pinecone CLI

Guide for AI agents and assistants contributing to this repository.

---

## Project overview

`pc` is the official Pinecone CLI — open source, public preview.

- **Entry point:** `cmd/pc/main.go` → `cliRootCmd.Execute()` (package `internal/pkg/cli/command/root`)
- **Module:** `github.com/pinecone-io/cli`

---

## Tech stack

| Concern | Library |
|---|---|
| CLI framework | `github.com/spf13/cobra` — hierarchical subcommands |
| Config | `github.com/spf13/viper` — YAML config files in `~/.config/pinecone/` |
| Pinecone client | `github.com/pinecone-io/go-pinecone/v5` |
| Auth | `github.com/golang-jwt/jwt/v5`, `golang.org/x/oauth2` |
| TUI / output | `github.com/charmbracelet/bubbletea`, `github.com/charmbracelet/lipgloss`, `github.com/fatih/color` |
| Logging | `github.com/rs/zerolog` |
| Testing | `github.com/stretchr/testify` |
| Go version | 1.24.0 |

---

## Architecture

### Command hierarchy

```
pc [group] [command] [subcommand]
# e.g. pc index vector upsert
```

Four root groups (defined in `internal/pkg/cli/command/root/root.go`):

| Group | Commands |
|---|---|
| **Auth** | auth, login, logout, target, whoami |
| **Admin** | organization, project, apiKey |
| **VectorDB** | index, collection, backup |
| **Misc** | version, config |

### Presenter system

Output formatting lives in `internal/pkg/utils/presenters/` — one file per output type. Command handlers stay thin; all rendering is delegated to presenters.

### Credential resolution

Three auth methods, resolved in `internal/pkg/utils/sdk/` and `internal/pkg/utils/oauth/`:

1. User login (OAuth2 PKCE flow)
2. Service account
3. API key

Credentials are stored in `~/.config/pinecone/`.

### Input handling

`internal/pkg/utils/argio/` unifies three JSON input methods:

| Method | Usage |
|---|---|
| Inline | Pass JSON directly as a flag value |
| File | Pass a `.json` or `.jsonl` file path |
| Stdin | Pass `-` to read from stdin |

### `--json` flag

Available on most commands. Forces structured, machine-readable output. Also activated automatically when stdout is not a TTY.

---

## Key directories

| Path | Purpose |
|---|---|
| `cmd/pc/` | Entry point |
| `internal/pkg/cli/command/` | All command implementations, organized by domain (`index/`, `auth/`, `project/`, etc.) |
| `internal/pkg/utils/` | Shared utilities: presenters, sdk, oauth, argio, flags, configuration, etc. |
| `test/e2e/` | End-to-end integration tests |

---

## Build & test commands

```bash
just build                                    # goreleaser build for current OS → ./dist/
just test-unit                                # go test -v ./... (no external deps required)
just test-e2e                                 # builds binary, runs E2E suite (requires credentials)
go test -v -run TestNameHere ./internal/...   # run a single test by name
```

E2E tests require a `.env` file at the repo root:

```
PINECONE_API_KEY=...
PINECONE_CLIENT_ID=...
PINECONE_CLIENT_SECRET=...
```

---

## Code style & conventions

- Always run `gofmt`. Zero `go vet` warnings allowed.
- `PascalCase` for exported symbols, `camelCase` for unexported. Acronyms uppercase: `URL`, `HTTP`, `API`, `ID`.
- All exported symbols require a GoDoc comment starting with the symbol name.
- Return errors; do not panic. Wrap with `fmt.Errorf("context: %w", err)`.
- Pass `context.Context` as the first parameter of any function that does I/O.
- Unit test names end with `Unit`. E2E tests live in `test/e2e/`.

---

## Scratch directory

`scratch/` is gitignored — use it freely for agent-generated temporary files, plans, and notes.
