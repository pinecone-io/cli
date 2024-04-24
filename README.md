# Pinecone CLI

`pinecone` is Pinecone on the command line. See the [Pinecone CLI PRD](https://www.notion.so/PRD-Pinecone-CLI-59fda5da83bc4e3a8593b74056914cd1?pm=c)

## Building the CLI

To build, run `make`. The built artifact will be placed into the `bin/` folder.

```
brew install goreleaser/tap/goreleaser
goreleaser build --single-target --snapshot --clean
```

For manual testing in development, you can run commands like this

```shell
./dist/pinecone_darwin_arm64/pinecone auth set-api-key "foo"
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