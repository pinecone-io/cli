/*
Copyright © 2025 Pinecone Systems, Inc.
*/
package main

import (
	"os"
	"path/filepath"

	cliRootCmd "github.com/pinecone-io/cli/internal/pkg/cli/command/root"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
)

func main() {
	executableName := filepath.Base(os.Args[0])
	if executableName == "pinecone" {
		pcio.Fprintf(os.Stderr, "⚠️  Warning: The '%s' command is deprecated. Please use '%s' instead.", style.Code("pinecone"), style.Code("pc"))
	}

	cliRootCmd.Execute()
}
