# AGENTS.md — Pinecone CLI

Guide for AI agents and assistants contributing to this repository.

---

## Project overview

`pc` is the official Pinecone CLI — open source.

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
| **VectorDB** | index |
| **Misc** | version, config |

`pc index` itself has internal subgroups:

| Subgroup | Subcommands |
|---|---|
| _(top-level index ops)_ | describe, list, create, configure, delete, describe-stats |
| **Data** | record, vector |
| **Namespace** | namespace |
| **Index Management** | backup, restore, collection, import |

### Presenter system

Output formatting lives in `internal/pkg/utils/presenters/` — one file per output type. Command handlers stay thin; all rendering is delegated to presenters.

### Credential resolution

Three auth methods, resolved in `internal/pkg/utils/sdk/` and `internal/pkg/utils/oauth/`:

| Priority | Method | Env vars |
|---|---|---|
| 1 (highest) | API key | `PINECONE_API_KEY` |
| 2 | Service account | `PINECONE_CLIENT_ID`, `PINECONE_CLIENT_SECRET` |
| 3 | User login (OAuth2 PKCE) | stored token in `~/.config/pinecone/secrets.yaml` |

API key mode does not require a target project to be set. OAuth and service account modes do.

Credentials and tokens are stored in `~/.config/pinecone/`. The directory holds two files:

- `config.yaml` — non-secret settings (`color`, `environment`)
- `secrets.yaml` — credentials and tokens (`api_key`, `client_id`, `client_secret`, `oauth2_token`); created with `0600` permissions

The `pc config` command manages these settings via a key-based interface. Valid keys: `api-key`, `color`, `environment`. Core subcommands: `get <key>`, `set <key> <value>`, `unset <key>`, `list`, `describe <key>`. Legacy shorthand commands (`set-api-key`, `set-color`, `set-environment`, `get-api-key`) still exist but the key-based interface is canonical.

### Input handling

`internal/pkg/utils/argio/` unifies three JSON input methods:

| Method | Usage |
|---|---|
| Inline | Pass JSON directly as a flag value |
| File | Pass a `.json` or `.jsonl` file path |
| Stdin | Pass `-` to read from stdin |

### `--json` flag

Available on most commands. Must be set explicitly to force structured, machine-readable output — it is **not** inferred from whether stdout is a TTY. (Color is suppressed automatically on a non-TTY, but the data output format is not.) The interactive auth flows — `login`, `target`, and `auth` — are the exception: they additionally infer JSON output when stdout is not a TTY, since an agent on a non-TTY can't drive their prompts.

---

## Patterns for writing commands

### Output: messages, styling, and exits

Never call `fmt.Print*` or `os.Exit` directly. Use the dedicated packages:

| Package | Use for |
|---|---|
| `internal/pkg/utils/msg/` | User-facing messages to stderr: `msg.SuccessMsg`, `msg.WarnMsg`, `msg.InfoMsg`, `msg.FailMsg`, `msg.FailJSON` |
| `internal/pkg/utils/style/` | Inline styled text: `style.Code()`, `style.URL()`, `style.Emphasis()`, `style.Spinner()` |
| `internal/pkg/utils/exit/` | Process exits: `exit.Success()`, `exit.Error()`, `exit.ErrorMsg()` |
| `internal/pkg/utils/text/` | JSON marshaling: `text.InlineJSON()`, `text.IndentJSON()` (never `json.Marshal` directly — it HTML-escapes) |

`msg.FailJSON(isJSON, ...)` writes an error to stderr in the right format (plain or JSON) based on whether `--json` is active. Always use it when a command can fail in JSON mode.

`style.*` functions respect the user's `color` setting and TTY state — never use raw ANSI codes or `lipgloss` directly in command handlers.

### Help text

Use `help.Long()` and `help.Examples()` to write multiline help strings. Both functions accept heredocs and normalize indentation:

```go
Long: help.Long(`
    Brief description here.

    More detail on second paragraph.
`),
Example: help.Examples(`
    pc index create --name my-index --dimension 1536
    pc index create --name my-index --help
`),
```

`help.Examples()` automatically prepends `$` to each command line and indents the block. Use the predefined group constants (`help.GROUP_AUTH`, `help.GROUP_VECTORDB`, `help.GROUP_INDEX_DATA`, etc.) when registering subcommand groups with `cmd.AddGroup()`.

### Custom flag types for JSON input

For flags that accept JSON arrays or objects, use the typed flag wrappers in `internal/pkg/utils/flags/` instead of plain string flags. They automatically support inline JSON, file paths, and stdin (`-`):

| Type | Use for |
|---|---|
| `flags.JSONObject` | A JSON object (`{...}`) |
| `flags.Float32List` | JSON array of float32 |
| `flags.Int32List` | JSON array of int32 |
| `flags.UInt32List` | JSON array of uint32 |
| `flags.StringList` | JSON array of strings |

Register with `cmd.Flags().Var(&myFlag, "flag-name", "description")`.

### Command testability: narrow service interfaces

Commands should not depend on the full SDK client directly. Define a narrow interface for the operations the command needs, and accept it as a parameter:

```go
type CreateIndexService interface {
    CreateIndex(ctx context.Context, req *pinecone.CreateIndexRequest) (*pinecone.Index, error)
}

func RunCreateIndex(ctx context.Context, svc CreateIndexService, opts CreateIndexOptions) error { ... }
```

This lets unit tests pass a mock without spinning up a real Pinecone client. See `internal/pkg/cli/command/index/create.go` for a reference implementation.

### Target org and project

Commands that operate on project-scoped resources must resolve the target org and project at runtime. Use:

```go
import "github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"

orgId, err := state.GetTargetOrgId()
projectId, err := state.GetTargetProjectId()
```

Both return an error if no target has been set (e.g., user hasn't run `pc target`). Surface the error through `exit.Error` so the user gets a clear prompt to run `pc target`.

### Logging

Use `internal/pkg/utils/log/` for internal diagnostics. Output goes to stderr and is gated by `PINECONE_LOG_LEVEL`:

```go
log.Debug(ctx).Str("index", name).Msg("fetching index")
log.Info(ctx).Msg("authenticated")
```

Never use `log/` for user-facing messages — that's `msg/`'s job.

---

## Notable environment variables

| Variable | Effect |
|---|---|
| `PINECONE_API_KEY` | Overrides stored API key (highest-priority auth) |
| `PINECONE_CLIENT_ID` / `PINECONE_CLIENT_SECRET` | Service account credentials |
| `PINECONE_ENVIRONMENT` | Target environment: `production` (default) or `staging` |
| `PINECONE_LOG_LEVEL` | Enable diagnostic logging: `INFO`, `DEBUG`, or `TRACE` |
| `PINECONE_CLI_MAX_JSON_BYTES` | Override the JSON input size limit (default 1 GiB) |

---

## Key directories

| Path | Purpose |
|---|---|
| `cmd/pc/` | Entry point |
| `internal/pkg/cli/command/` | All command implementations, organized by domain (`index/`, `auth/`, `project/`, etc.) |
| `internal/pkg/utils/` | Shared utilities: presenters, sdk, oauth, argio, flags, configuration, etc. |
| `scripts/` | Shell helpers: `install.sh`, `uninstall.sh` |
| `test/e2e/` | End-to-end integration tests |

---

## Build & test commands

```bash
just build                                    # goreleaser build for current OS → ./dist/
just build-all                                # goreleaser build for all supported OSes
just test-unit                                # go test -v ./... (no external deps required)
just test-e2e                                 # builds binary, runs E2E suite (requires credentials)
just gen-manpages                             # generate man pages → ./man/
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
- Unit test function names follow the pattern `Test<FunctionName>` or `Test<FunctionName>_<Scenario>`. E2E tests live in `test/e2e/` and use the `//go:build e2e` build tag — `just test-unit` excludes them automatically.

---

## Scratch directory

`scratch/` is gitignored — use it freely for agent-generated temporary files, plans, and notes.
