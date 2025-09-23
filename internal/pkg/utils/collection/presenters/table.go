package presenters

import (
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/go-pinecone/v4/pinecone"
)

// PrintCollectionTable creates and renders a table for collection list with proper column formatting
func PrintCollectionTable(collections []*pinecone.Collection) {
	if len(collections) == 0 {
		fmt.Println("No collections found.")
		return
	}

	// Define table columns
	columns := []presenters.Column{
		{Title: "NAME", Width: 20},
		{Title: "DIMENSION", Width: 10},
		{Title: "SIZE", Width: 8},
		{Title: "STATUS", Width: 12},
		{Title: "VECTORS", Width: 10},
		{Title: "ENVIRONMENT", Width: 15},
	}

	// Convert collections to table rows
	rows := make([]presenters.Row, len(collections))
	for i, collection := range collections {
		rows[i] = presenters.Row{
			collection.Name,
			fmt.Sprintf("%d", collection.Dimension),
			presenters.FormatSize(int(collection.Size)),
			string(collection.Status), // Use plain string instead of colorized version
			fmt.Sprintf("%d", collection.VectorCount),
			collection.Environment,
		}
	}

	presenters.PrintTable(presenters.TableOptions{
		Columns: columns,
		Rows:    rows,
	})

	fmt.Println()
	hint := fmt.Sprintf("Use %s to see collection details", style.Code("pc collection describe <name>"))
	fmt.Println(style.Hint(hint))
}
