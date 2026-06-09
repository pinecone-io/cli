# Contributing to the Pinecone CLI

## Building the CLI

1. [Install Go](https://go.dev/doc/install) if you do not have it already.

2. [Install `just`](https://github.com/casey/just?tab=readme-ov-file#installation) to run the recipes in the [justfile](justfile).

3. Install goreleaser:

```bash
brew install --cask goreleaser/tap/goreleaser
```

4. Clone the repo and build:

```bash
git clone git@github.com:pinecone-io/cli.git
just build
```

This builds a binary for your current OS and places it under `./dist/`. The exact subdirectory depends on your platform (e.g. `pc_darwin_arm64_v8.0/pc` on Apple Silicon). Run `ls dist/` to find the right path.

For a quick alias during development:

```bash
alias pc-dev="$(ls -d dist/pc_darwin_arm64*/pc dist/pc_darwin_all/pc 2>/dev/null | head -1)"
```

## Verifying your build

```bash
# See help
./dist/pc_darwin_arm64_v8.0/pc --help

# Authenticate
./dist/pc_darwin_arm64_v8.0/pc auth login
# or set an API key directly
./dist/pc_darwin_arm64_v8.0/pc config set api-key "YOUR_API_KEY"

# Spot-check index operations
./dist/pc_darwin_arm64_v8.0/pc index list
./dist/pc_darwin_arm64_v8.0/pc index create --name my-index --dimension 1536 --metric cosine --cloud aws --region us-east-1
./dist/pc_darwin_arm64_v8.0/pc index describe --index-name my-index
./dist/pc_darwin_arm64_v8.0/pc index delete --index-name my-index
```

For full usage documentation, see the [README](./README.md).

## Running tests

```bash
# Unit tests — no external dependencies required
just test-unit

# Run a single test by name
go test -v -run TestNameHere ./internal/...

# E2E tests — builds the binary and runs against real Pinecone APIs
# Requires credentials to be set in the environment or an .env file
just test-e2e
```

E2E tests require the following credentials to be available, either as environment variables or in a `.env` file at the repo root:

```
PINECONE_API_KEY=...
PINECONE_CLIENT_ID=...
PINECONE_CLIENT_SECRET=...
```

E2E tests should use the `//go:build e2e` build tag, so they are excluded from `just test-unit` automatically.

## Troubleshooting

- Configuration files are stored in `~/.config/pinecone/`.
- Enable debug output with `PINECONE_LOG_LEVEL=DEBUG`.
- Check which environment you're pointed at: `cat ~/.config/pinecone/config.yaml` or `pc config get environment`. If things aren't working as expected, confirm the `environment` setting is correct (`production` or `staging`). To switch: `pc config set environment production`.

## Making a Pull Request

Fork this repo and open a PR with your changes. Before submitting:

- Run `gofmt` on all changed files.
- Run `go vet ./...` and resolve any warnings.
- Run `just test-unit` and ensure tests pass.

## Releasing the CLI

To make a new release, tag a commit and push to remote. CI handles the rest.

```sh
# Ensure main is clean and up to date
git checkout main
git pull
git status

# Confirm the tip of main builds
goreleaser build --clean --snapshot

# Check existing tags to pick the next version
git tag --list

# Tag and push (must start with "v" to trigger CI)
git tag v0.1.0
git push --tags
```

The [publish workflow](https://github.com/pinecone-io/cli/actions/workflows/publish.yaml) uses [goreleaser](https://goreleaser.com/) to build binaries for all supported platforms, publish artifacts to the GitHub Releases page, and update the Homebrew tap. The `.goreleaser.yaml` file is the authoritative configuration if anything needs adjusting.

Within a few minutes of pushing a tag you should see:

- A new entry on the [Releases page](https://github.com/pinecone-io/cli/releases) with built artifacts attached.
- An automatic update to the [Homebrew tap](https://github.com/pinecone-io/homebrew-tap).

Users on Homebrew can upgrade with:

```sh
brew update
brew upgrade --cask pinecone
```

Users who installed via the install script can upgrade by re-running it:

```sh
curl -fsSL https://pinecone.io/install.sh | sh
```
