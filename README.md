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
2. **Download the appropriate binary** for your platform and architecture
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

- macOS (x86_64, ARM64)
- Linux (x86_64, ARM64)
- Windows (x86_64)

### Build from source

To learn about the steps involved in building from source, see [CONTRIBUTING](./CONTRIBUTING.md)

## Authentication

There are three ways to authenticate the Pinecone CLI: through a web browser with user login, using a service account, or with an API key.

This table describes the Pinecone operations supported by each authentication method:

| Method          | Admin API | Control plane |
| :-------------- | :-------- | :------------ |
| User login      | ✅        | ✅            |
| Service account | ✅        | ✅            |
| API key         | ❌        | ✅            |

- Admin API–related commands (organization and project management, API key operations):

  - `pc organization` (`list`, `describe`, `update`, `delete`)
  - `pc project` (`create`, `list`, `describe`, `update`, `delete`)
  - `pc api-key` (`create`, `list`, `describe`, `update`, `delete`)

- Control plane–related commands (index management):
  - `pc index` (`create`, `list`, `describe`, `configure`, `delete`)

### 1. User Login (Recommended for Interactive Use)

The standard authentication method for interactive use. Provides full access to the Admin API and control/data plane operations. You will have access to all organizations associated with the account.

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

### 2. Service Account Authentication

Use [service account](https://docs.pinecone.io/guides/organizations/manage-service-accounts) client credentials for authentication. Service accounts are scoped to a single organization, but you can manage projects and set a target project context.

```bash
# Will prompt you to pick a target project from the projects available to the service account
pc auth configure --client-id "YOUR_CLIENT_ID" --client-secret "YOUR_CLIENT_SECRET"

# Specify a target project when configuring the service account
pc auth configure --client-id "client-id" --client-secret "client-secret" --project-id "project-id"
```

### 3. API Key Authentication

Use a project API key directly. Provides access to control/data plane operations only (no Admin API access). If an API key is set directly, it will override any configured target organization and project context.

```bash
pc auth configure --api-key "YOUR_API_KEY"

# alternatively
pc config set-api-key "YOUR_API_KEY"
```

For more detailed information, see the [CLI Authentication documentation](https://docs.pinecone.io/reference/cli/authentication).

## Quick Start

After installation, get started with authentication:

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
