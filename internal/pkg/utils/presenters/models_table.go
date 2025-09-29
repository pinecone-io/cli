package presenters

import (
	"fmt"
	"strconv"

	"github.com/pinecone-io/go-pinecone/v4/pinecone"
)

// PrintModelsTable creates and renders a table showing model information
func PrintModelsTable(models []pinecone.ModelInfo) {
	if len(models) == 0 {
		fmt.Println("No models found.")
		return
	}

	// Define table columns
	columns := []Column{
		{Title: "Model", Width: 25},
		{Title: "Type", Width: 8},
		{Title: "Vector Type", Width: 12},
		{Title: "Dimension", Width: 10},
		{Title: "Provider", Width: 15},
		{Title: "Description", Width: 40},
	}

	// Convert models to table rows
	rows := make([]Row, len(models))
	for i, model := range models {
		dimension := "-"
		if model.DefaultDimension != nil {
			dimension = strconv.Itoa(int(*model.DefaultDimension))
		}

		vectorType := "-"
		if model.VectorType != nil {
			vectorType = *model.VectorType
		}

		provider := "-"
		if model.ProviderName != nil {
			provider = *model.ProviderName
		}

		// Truncate description if too long
		description := model.ShortDescription
		if len(description) > 35 {
			description = description[:32] + "..."
		}

		rows[i] = Row{
			model.Model,
			model.Type,
			vectorType,
			dimension,
			provider,
			description,
		}
	}

	// Print the table
	PrintTable(TableOptions{
		Columns: columns,
		Rows:    rows,
	})
}

// PrintModelsTableWithTitle creates and renders a models table with a title
func PrintModelsTableWithTitle(title string, models []pinecone.ModelInfo) {
	fmt.Println()
	fmt.Printf("%s\n\n", title)
	PrintModelsTable(models)
	fmt.Println()
}
