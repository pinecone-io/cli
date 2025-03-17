# Pinecone CLI

> ⚠️ **Warning:** This SDK is still in an alpha state. While it is mostly built out and functional, it may undergo changes as we continue to improve the UX. Please try it out and give us your feedback, but be aware that updates may introduce breaking changes.

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

In order to use the Pinecone CLI you will need to authenticate with Pinecone services. This can be done either with an API key, or using the `pinecone login` flow to authenticate with a Pinecone account via your browser.

```shell
pinecone --help

# If you have PINECONE_API_KEY set in your environment you can begin working with the CLI
pinecone index list

# To set an API key manually, you can use the config command
pinecone config set-api-key "YOUR_API_KEY"

# Additionally, you can authenticate through the browser using the login command
pinecone login

# To clear your current login state or configured API key, you can use the logout command
pinecone logout
```

If an API key is configured along with using `pinecone login`, the CLI will default to using the API key over the authentication token.

If there has been an API key set using `pinecone config set-api-key`, and `PINECONE_API_KEY` is also present in the environment, the API set in the CLI config will be used over the environment key.

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
