# Pinecone CLI

`pinecone` is Pinecone on the command line. See the [Pinecone CLI PRD](https://www.notion.so/PRD-Pinecone-CLI-59fda5da83bc4e3a8593b74056914cd1?pm=c)

## Building the CLI

1. [Install golang](https://go.dev/doc/install) if you do not have it already

2. Install goreleaser

```
brew install goreleaser/tap/goreleaser
```

3. Build the CLI

```
goreleaser build --single-target --snapshot --clean
```

For manual testing in development, you can run commands like this

```shell
./dist/pinecone_darwin_arm64/pinecone login
./dist/pinecone_darwin_arm64/pinecone index list
# etc
```

## Usage

```shell
# See help
./dist/pinecone_darwin_arm64/pinecone --help

# Set authorization credentials - set an API key directly, or log in via the OAuth flow
./dist/pinecone_darwin_arm64/pinecone config set-api-key
./dist/pinecone_darwin_arm64/pinecone login

# Check currently configured API key
./dist/pinecone_darwin_arm64/pinecone config get-api-key

# Do index operations
./dist/pinecone_darwin_arm64/pinecone index --help

# Create serverless indexes.
./dist/pinecone_darwin_arm64/pinecone index create-serverless --help
./dist/pinecone_darwin_arm64/pinecone index create-serverless --name example-index --dimension 1536 --metric cosine --cloud aws --region us-west-2
./dist/pinecone_darwin_arm64/pinecone index create-serverless --name="example-index" --dimension=1536 --metric="cosine" --cloud="aws" --region="us-west-2"
./dist/pinecone_darwin_arm64/pinecone index create-serverless -n example-index -d 1536 -m cosine -c aws -r us-west-2

# Describe index
./dist/pinecone_darwin_arm64/pinecone index describe --name "example-index"
./dist/pinecone_darwin_arm64/pinecone index describe --name "example-index" --json

# List indexes
./dist/pinecone_darwin_arm64/pinecone index list
./dist/pinecone_darwin_arm64/pinecone index list --json

# Delete index
./dist/pinecone_darwin_arm64/pinecone index delete --name "example-index"
```

## Releasing the CLI

To make a new release, you simply tag a commit with a version and push it. The heavy lifting all happens in CI.

Something along these lines:

```
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

To consume the update from Homebrew, users should run

```sh
brew update
brew upgrade pinecone
```