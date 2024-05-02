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

# Set credentials (proper login will come later)
./dist/pinecone_darwin_arm64/pinecone auth set-api-key

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