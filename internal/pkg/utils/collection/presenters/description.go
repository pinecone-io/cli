package presenters

import (
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/go-pinecone/v4/pinecone"
)

func PrintDescribeCollectionTable(coll *pinecone.Collection) {
	log.Debug().Str("name", coll.Name).Msg("Printing collection description")

	// Build rows for the table
	var rows []presenters.Row

	// Essential information
	rows = append(rows, presenters.Row{"Name", coll.Name})
	rows = append(rows, presenters.Row{"State", ColorizeCollectionStatus(coll.Status)})
	rows = append(rows, presenters.Row{"Dimension", fmt.Sprintf("%d", coll.Dimension)})
	rows = append(rows, presenters.Row{"Size", presenters.FormatSize(int(coll.Size))})
	rows = append(rows, presenters.Row{"Vector Count", fmt.Sprintf("%d", coll.VectorCount)})
	rows = append(rows, presenters.Row{"Environment", coll.Environment})

	// Print each row with right-aligned first column and secondary text styling
	for _, row := range rows {
		if len(row) >= 2 {
			// Right align the first column content
			rightAlignedFirstCol := fmt.Sprintf("%20s", row[0])

			// Apply secondary text styling to the first column
			styledFirstCol := style.SecondaryTextStyle().Render(rightAlignedFirstCol)

			// Print the row
			rowText := fmt.Sprintf("%s  %s", styledFirstCol, row[1])
			fmt.Println(rowText)
		}
	}
	// Add spacing after the last row
	fmt.Println()
}

func ColorizeCollectionStatus(state pinecone.CollectionStatus) string {
	switch state {
	case pinecone.CollectionStatusReady:
		return style.SuccessStyle().Render(string(state))
	case pinecone.CollectionStatusInitializing, pinecone.CollectionStatusTerminating:
		return style.WarningStyle().Render(string(state))
	}

	return string(state)
}
