/*
Copyright © 2025 Pinecone Systems, Inc.
*/
package main

import (
	"fmt"
	"os"
	"path/filepath"

	cliRootCmd "github.com/pinecone-io/cli/internal/pkg/cli/command/root"
)

func main() {
	executableName := filepath.Base(os.Args[0])
	if executableName == "pinecone" {
		fmt.Fprintln(os.Stderr, "⚠️  Warning: The 'pinecone' command is deprecated. Please use 'pc' instead.")
	}

	cliRootCmd.Execute()
}
