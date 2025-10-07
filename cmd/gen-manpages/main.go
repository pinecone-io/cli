package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pinecone-io/cli/internal/pkg/cli/command/root"
	"github.com/pinecone-io/cli/internal/pkg/utils/manpage"
)

func main() {
	var (
		output  = flag.String("output", "./man/man1", "Output directory for the generated manpages")
		verbose = flag.Bool("verbose", false, "Verbose output when generating")
		help    = flag.Bool("help", false, "Show help for this command")
	)
	flag.Parse()

	// Display help
	if *help {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
		os.Exit(0)
	}

	if *verbose {
		fmt.Printf("Generating man pages for Pinecone CLI...\n")
		fmt.Printf("Output directory: %s\n", *output)
	}

	// Get root command for the CLI
	rootCmd := root.GetRootCmd()
	if rootCmd == nil {
		fmt.Fprintf(os.Stderr, "Error: Could not access root command\n")
		os.Exit(1)
	}

	// Ensure output directory exists
	if err := os.MkdirAll(*output, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Could not create output directory: %s\n", err)
		os.Exit(1)
	}

	// Generate man pages
	if err := manpage.GenerateManPages(rootCmd, *output); err != nil {
		fmt.Fprintf(os.Stderr, "Error generating man pages: %s\n", err)
		os.Exit(1)
	}

	// List generated files for verbose
	if *verbose {
		files, err := filepath.Glob(filepath.Join(*output, "*.1"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to list generated files: %s\n", err)
		} else if len(files) == 0 {
			fmt.Printf("No man pages found in output directory: %s\n", *output)
		} else {
			fmt.Printf("Generated %d man pages:\n", len(files))
			for _, file := range files {
				fmt.Printf(" - %s\n", filepath.Base(file))
			}
		}
	}

	fmt.Println("Pinecone CLI man pages generated successfully")
}
