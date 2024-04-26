package presenters

import (
	"fmt"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
)

func PrintTargetContext(context *state.TargetContext) {
	writer := NewTabWriter()

	columns := []string{"ATTRIBUTE", "VALUE"}
	header := strings.Join(columns, "\t") + "\n"
	fmt.Fprint(writer, header)

	fmt.Fprintf(writer, "Api\t%s\n", style.Emphasis(context.Api))
	fmt.Fprintf(writer, "Org\t%s\n", context.Org)
	fmt.Fprintf(writer, "Project\t%s\n", context.Project)

	writer.Flush()
}
