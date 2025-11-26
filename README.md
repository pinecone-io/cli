# Pinecone CLI

`pc` is Pinecone on the command line.

> ⚠️ **Note:** This CLI is in [public preview](https://docs.pinecone.io/assistant-release-notes/feature-availability) and does not yet support all features available through the Pinecone API. Please try it out and let us know of any feedback. You'll want to upgrade often as we address feedback and add additional features.

## Installation

### Homebrew (macOS, Linux)

The most convenient way to install the CLI on macOS and Linux is via [Homebrew](https://brew.sh).

If you don't have Homebrew installed, install it first:

```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

1. **Add the Pinecone tap** to your Homebrew configuration:

```bash
brew tap pinecone-io/tap
```

2. **Install the Pinecone CLI**:

```bash
brew install pinecone-io/tap/pinecone
```

3. **Verify the installation**:

```bash
pc --help
```

#### What is a Homebrew tap?

A [Homebrew tap](https://docs.brew.sh/Taps) is a third-party repository of Homebrew formulas. Our official tap at [pinecone-io/homebrew-tap](https://github.com/pinecone-io/homebrew-tap) contains the formula needed to install the Pinecone CLI.

#### Upgrading

To upgrade to the latest version:

```bash
brew update
brew upgrade pinecone
```

#### Uninstalling

To remove the CLI:

```bash
brew uninstall pinecone
```

To remove the Pinecone tap entirely:

```bash
brew untap pinecone-io/tap
```

### Download artifacts from release page (Linux, Windows, macOS)

For users who prefer not to use Homebrew or need specific platform binaries, we provide pre-built binaries for many platforms.

1. **Visit the [Releases page](https://github.com/pinecone-io/cli/releases)**
2. **Download the appropriate binary** for your platform and architecture from the "Assets" section.
3. **Make the binary executable** (Linux/macOS):
   ```bash
   chmod +x pc
   ```
4. **Move to a directory in your PATH** (optional but recommended):
   ```bash
   sudo mv pc /usr/local/bin/  # Linux/macOS
   # or on Windows, add the directory to your PATH
   ```
5. **Verify the installation**:
   ```bash
   pc --help
   ```

#### Supported platforms

- macOS: Intel (x86_64) and Apple Silicon (ARM64)
- Linux: x86_64, ARM64, and i386 architectures
- Windows: x86_64 and i386 architectures

### Build from source

To learn about the steps involved in building from source, see [CONTRIBUTING](./CONTRIBUTING.md)

## Authentication

There are three ways to authenticate the Pinecone CLI: through a web browser with user login, using a service account, or with an API key.

This table describes the Pinecone operations supported by each authentication method:

| Method          | Admin API | Control plane | Data plane |
| :-------------- | :-------- | :------------ | :--------- |
| User login      | ✅        | ✅            | ✅         |
| Service account | ✅        | ✅            | ✅         |
| API key         | ❌        | ✅            | ✅         |

- Admin API–related commands (organization and project management, API key operations):

  - `pc organization` (`list`, `describe`, `update`, `delete`)
  - `pc project` (`create`, `list`, `describe`, `update`, `delete`)
  - `pc api-key` (`create`, `list`, `describe`, `update`, `delete`)

- Control plane–related commands (index management):

  - `pc index` (`create`, `list`, `describe`, `configure`, `delete`, `describe-stats`)

- Data plane-related commands (index data management):

  - `pc index vector` (`upsert`, `query`, `fetch`, `list`, `update`, `delete`)

### 1. User Login (Recommended for Interactive use)

The standard authentication method for interactive use. Provides full access to the Admin API and control/data plane operations. When authenticated this way, you have access to all organizations associated with the account.

```bash
pc auth login
```

This command:

- Opens your browser to the Pinecone login page
- Automatically sets a target organization and project context
- Grants access to manage organizations, projects, and other account-level resources

**View and change your current target:**

```bash
pc target -s
pc target -o "ORGANIZATION_NAME" -p "PROJECT_NAME"
```

### 2. Service account authentication

Use [service account](https://docs.pinecone.io/guides/organizations/manage-service-accounts) client credentials for authentication. Service accounts are scoped to a single organization, but you can manage projects and set a target project context.

```bash
# Prompts you to pick a target project from the projects available to the service account
pc auth configure --client-id "YOUR_CLIENT_ID" --client-secret "YOUR_CLIENT_SECRET"

# Specify a target project when configuring the service account
pc auth configure --client-id "client-id" --client-secret "client-secret" --project-id "project-id"
```

### 3. API key authentication

Use a project API key directly. Provides access to control/data plane operations only (no Admin API access). If an API key is set directly, it overrides any configured target organization and project context.

```bash
pc auth configure --api-key "YOUR_API_KEY"

# alternatively
pc config set-api-key "YOUR_API_KEY"
```

For more detailed information, see the [CLI authentication](https://docs.pinecone.io/reference/cli/authentication) documentation.

## Data plane commands overview

Work with your vector data inside an index. These commands require `--index-name` and optionally `--namespace`:

- Ingest and manage records:
  - `pc index vector upsert` — insert or update vectors from JSON/JSONL
  - `pc index vector list` — list vectors (with pagination)
  - `pc index vector fetch` — fetch by IDs or metadata filter
  - `pc index vector update` — update a vector by ID or update many via metadata filter
  - `pc index vector delete` — delete by IDs, by filter, or delete all in a namespace
  - `pc index vector query` — nearest-neighbor search by values or vector ID
- Index statistics:
  - `pc index describe-stats` — show dimension, vector counts, namespace summary, and metadata field counts (optionally filtered)

Tip: add `--json` to many commands to get structured output.

## Quickstart

After installing the CLI, authenticate with user login or set an API key, verify your auth status, and list indexes associated with your automatically targeted project.

```bash
# Option 1: Login via browser (recommended)
pc auth login

# Option 2: Set API key directly
pc config set-api-key "YOUR_API_KEY"

# Verify authentication
pc auth whoami
pc auth status

# List your indexes
pc index list
```

### JSON input formats

Many flags accept JSON in three forms:

- Inline JSON for smaller payloads, for example:
  ```bash
  pc index vector fetch --index-name my-index --namespace demo --ids '["vec-1","vec-2"]'
  ```
- From a file ending in `.json` or `.jsonl`:
  ```bash
  pc index vector upsert --index-name my-index --namespace demo --body ./vectors.jsonl
  ```
- From stdin with `-`:
  ```bash
  cat vectors.jsonl | pc index vector upsert --index-name my-index --namespace demo --body -
  ```

Stdin can only be used with one flag at a time.

## Data plane quickstart (end-to-end)

### Body JSON schemas

For commands that accept a `--body` JSON payload, the CLI uses these schemas:

- UpsertBody — vectors of `pinecone.Vector` (see `https://pkg.go.dev/github.com/pinecone-io/go-pinecone/v5/pinecone#Vector`)
- QueryBody — fields: id, vector, `sparse_values` (see `https://pkg.go.dev/github.com/pinecone-io/go-pinecone/v5/pinecone#SparseValues`), filter, top_k, include_values, include_metadata
- FetchBody — fields: ids, filter, limit, pagination_token
- UpdateBody — fields: id, values, `sparse_values` (see `https://pkg.go.dev/github.com/pinecone-io/go-pinecone/v5/pinecone#SparseValues`), metadata, filter, dry_run

The following walkthrough creates an index, ingests vectors, and runs queries.

Prepare sample vectors (JSONL)

Create a file named `vectors.jsonl` with two lines:

```json
{"id":"vec-1","values":[0.1,0.2,0.3],"metadata":{"genre":"sci-fi","title":"Voyager"}}
{"id":"vec-2","values":[0.3,0.1,0.2],"metadata":{"genre":"fantasy","title":"Dragon"}}
```

Alternatively, you can upsert using a JSON object with a `vectors` array:

```json
{
  "vectors": [
    {
      "id": "vec-1",
      "values": [0.1, 0.2, 0.3],
      "metadata": { "genre": "sci-fi", "title": "Voyager" }
    },
    {
      "id": "vec-2",
      "values": [0.3, 0.1, 0.2],
      "metadata": { "genre": "fantasy", "title": "Dragon" }
    }
  ]
}
```

```bash
# Create a serverless index
pc index create --name my-index --dimension 3 --metric cosine --cloud aws --region us-east-1

# Upsert vectors into the index via JSON or JSONL
pc index vector upsert --index-name my-index --namespace my-namespace --body ./vectors.jsonl

# List vectors
pc index vector list --index-name my-index --namespace my-namespace

# Fetch a vector by ID
pc index vector fetch --index-name my-index --namespace my-namespace --ids '["vec-1"]'

# Query by dense vector values
pc index vector query --index-name my-index --namespace my-namespace --vector '[0.1,0.2,0.3]' --top-k 3 --include-metadata

# Query by existing vector ID
pc index vector query --index-name my-index --namespace my-namespace --id vec-1 --top-k 3
```
