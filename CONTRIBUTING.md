# Pinecone CLI

`pc` is Pinecone on the command line. See the [Pinecone CLI PRD](https://www.notion.so/PRD-Pinecone-CLI-59fda5da83bc4e3a8593b74056914cd1?pm=c)

## Building the CLI

1. [Install golang](https://go.dev/doc/install) if you do not have it already

2. Install goreleaser

```
brew install --cask goreleaser/tap/goreleaser
```

3. Build the CLI

```
goreleaser build --single-target --snapshot --clean
```

For manual testing in development, you can run commands like this

```shell
./dist/pc_darwin_arm64/pc login
./dist/pc_darwin_arm64/pc index list
# etc
```

## Usage

```shell
# See help
./dist/pc_darwin_arm64/pc --help

# Set authorization credentials - set an API key directly, or log in via the OAuth flow
./dist/pc_darwin_arm64/pc config set-api-key
./dist/pc_darwin_arm64/pc login

# Check currently configured API key
./dist/pc_darwin_arm64/pc config get-api-key

# Do index operations
./dist/pc_darwin_arm64/pc index --help

# Create serverless indexes.
./dist/pc_darwin_arm64/pc index create-serverless --help
./dist/pc_darwin_arm64/pc index create-serverless --name example-index --dimension 1536 --metric cosine --cloud aws --region us-west-2
./dist/pc_darwin_arm64/pc index create-serverless --name="example-index" --dimension=1536 --metric="cosine" --cloud="aws" --region="us-west-2"
./dist/pc_darwin_arm64/pc index create-serverless -n example-index -d 1536 -m cosine -c aws -r us-west-2

# Describe index
./dist/pc_darwin_arm64/pc index describe --name "example-index"
./dist/pc_darwin_arm64/pc index describe --name "example-index" --json

# List indexes
./dist/pc_darwin_arm64/pc index list
./dist/pc_darwin_arm64/pc index list --json

# Delete index
./dist/pc_darwin_arm64/pc index delete --name "example-index"
```

## Troubleshooting

Some facts that could be useful:

- Configuration files are stored in `~/.config/pinecone`
- You can enable debug output with the `PINECONE_LOG_LEVEL=DEBUG` env var
- Are you pointed at the correct environment? The current value of the environment setting (i.e. prod or staging) is controlled through `pc config set-environment staging` is not clearly surfaced through the printed output. If things aren't working as you expect, you might be pointed in the wrong place. See `cat ~/.config/pinecone/config.yaml` to confirm.

## Development Practices & Tools

This project follows several established patterns and provides utilities to ensure consistency across the codebase.

### Output Functions & Quiet Mode

The CLI supports a `-q` (quiet) flag that suppresses non-essential output while preserving essential data. Follow these guidelines:

**Use `pcio` functions for:**

- User-facing messages (success, error, warning, info)
- Progress indicators and status updates
- Interactive prompts and confirmations
- Help text and documentation
- Any output that should be suppressed with `-q` flag

**Use `fmt` functions for:**

- Data output from informational commands (list, describe)
- JSON output that should always be displayed
- Table rendering and structured data display
- Any output that should NOT be suppressed with `-q` flag

```go
// ✅ Correct usage
pcio.Println("Creating index...")  // User message - suppressed with -q
msg.SuccessMsg("Index created!")   // User message - suppressed with -q
fmt.Println(jsonData)              // Data output - always displayed

// ❌ Incorrect usage
pcio.Println(jsonData)             // Wrong! Data would be suppressed
fmt.Println("Creating index...")   // Wrong! Ignores quiet mode
```

### Error Handling

Use the centralized error handling utilities:

```go
// For API errors with structured responses
errorutil.HandleIndexAPIError(err, cmd, args)

// For program termination
exit.Error(err)        // Logs error and exits with code 1
exit.ErrorMsg("msg")   // Logs message and exits with code 1
exit.Success()         // Logs success and exits with code 0
```

### User Messages & Styling

Use the `msg` package for consistent user messaging:

```go
msg.SuccessMsg("Operation completed successfully!")
msg.FailMsg("Operation failed: %s", err)
msg.WarnMsg("This will delete the resource")
msg.InfoMsg("Processing...")
msg.HintMsg("Use --help for more options")

// Multi-line messages
msg.WarnMsgMultiLine("Warning 1", "Warning 2", "Warning 3")
```

Use the `style` package for consistent text formatting:

```go
style.Heading("Section Title")
style.Emphasis("important text")
style.Code("command-name")
style.URL("https://example.com")
```

### Interactive Components

For user confirmations, use the interactive package:

```go
result := interactive.AskForConfirmation("Delete this resource?")
switch result {
case interactive.ConfirmationYes:
    // Proceed with deletion
case interactive.ConfirmationNo:
    // Cancel operation
case interactive.ConfirmationQuit:
    // Exit program
}
```

### Table Rendering

Use the `presenters` package for consistent table output:

```go
// For data tables (always displayed, not suppressed by -q)
presenters.PrintTable(presenters.TableOptions{
    Columns: []presenters.Column{{Title: "Name", Width: 20}},
    Rows:    []presenters.Row{{"example"}},
})

// For index-specific tables
presenters.PrintIndexTableWithIndexAttributesGroups(indexes, groups)
```

### Testing Utilities

Use the `testutils` package for consistent command testing:

```go
// Test command arguments and flags
tests := []testutils.CommandTestConfig{
    {
        Name:         "valid arguments",
        Args:         []string{"my-arg"},
        Flags:        map[string]string{"json": "true"},
        ExpectError:  false,
        ExpectedArgs: []string{"my-arg"},
    },
}
testutils.TestCommandArgsAndFlags(t, cmd, tests)

// Test JSON flag configuration
testutils.AssertJSONFlag(t, cmd)
```

### Validation Utilities

Use centralized validation functions:

```go
// For index name validation
index.ValidateIndexNameArgs(cmd, args)

// For other validations, check the respective utility packages
```

### Logging

Use structured logging with the `log` package:

```go
log.Debug().Str("index", name).Msg("Creating index")
log.Error().Err(err).Msg("Failed to create index")
log.Info().Msg("Operation completed")
```

### Configuration Management

Use the configuration utilities for consistent config handling:

```go
// Get current state
org := state.TargetOrg.Get()
proj := state.TargetProj.Get()

// Configuration files are managed through the config package
```

## Making a Pull Request

Please fork this repo and make a PR with your changes. Run `gofmt` and `goimports` on all proposed
code changes. Code that does not adhere to these formatters will not be merged.

## Releasing the CLI

To make a new release, you simply tag a commit with a version and push it. The heavy lifting all happens in CI.

Something along these lines:

```sh
# Pull and ensure you have no uncomitted changes
git checkout main
git pull
git status

# Ensure the tip of main actually builds
gorelaser build --clean --snapshot

# Look at what version tags have previously been used
git tag --list

# Based on the previous history and the nature of the
# new stuff in the code you are releasing, choose a
# tag that makes sense for the next version.
#
# The tag must start with "v" to trigger the CI stuff.
git tag v0.0.40

# Push the tag to github
git push --tags
```

From there, everything happens in this [publish workflow](https://github.com/pinecone-io/cli/actions/workflows/publish.yaml) which is using [goreleaser](https://goreleaser.com/) to handle the process of building binaries for different platforms, packing them into archives, publishing those artifacts on github, and updating our homebrew formula so those updates are easily installable on mac. In the future this will probably expand to cover more forms of distribution. If anything breaks down in this process, the `.goreleaser.yaml` file is probably where your attention will be needed but so far it has been very reliable.

Within a few minutes of pushing tags, you should see:

- A new update to the [Releases page](https://github.com/pinecone-io/cli/releases) with built artifacts attached. If you want to be fancy, you can edit the text there to give a more narrative overview of what is in the release. But for these early iterations we're just pushing and shipping without a lot of ceremony.
- Updates to to the [Homebrew tap](https://github.com/pinecone-io/homebrew-tap) should happen automatically

To consume the update from Homebrew (assuming they have previously installed it from homebrew), users should run

```sh
brew update
brew upgrade pinecone
```
