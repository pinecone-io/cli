# Pinecone CLI

`pinecone` is Pinecone on the command line. 

This CLI is still in an alpha state and does not support every operation available through our REST apis. Please try it out and give us your feedback, but also be prepared to upgrade as we continue building out the feature set and improving the UX.

## Installing

### Via Homebrew (Mac)

The most convenient way to install this is via [Homebrew](https://brew.sh)

```brew
brew tap pinecone-io/tap
brew install pinecone-io/tap/pinecone

pinecone --help
```

If you have previously installed and would like to upgrade to the latest version, run

```
brew update
brew upgrade pinecone
```

### Download artifacts from release page (Linux, Windows)

We have pre-built binaries for many platforms available on the [Releases](https://github.com/pinecone-io/cli/releases) page.

### Build from source

To learn about the steps involved in building from source, see [CONTRIBUTING](./CONTRIBUTING.md)

## Usage

```shell
pinecone --help

pinecone login
```

### Managing indexes

```sh
# Learn about supported index operations
pinecone index --help

# Create serverless indexes.
pinecone index create-serverless --help
pinecone index create-serverless --name example-index --dimension 1536 --metric cosine --cloud aws --region us-west-2
pinecone index create-serverless --name="example-index" --dimension=1536 --metric="cosine" --cloud="aws" --region="us-west-2"
pinecone index create-serverless -n example-index -d 1536 -m cosine -c aws -r us-west-2

# Describe index
pinecone index describe --name "example-index"
pinecone index describe --name "example-index" --json

# List indexes
pinecone index list
pinecone index list --json

# Delete index
pinecone index delete --name "example-index"
```
