package presenters

import (
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
)

func labelUnsetIfEmpty(value string) string {
	if value == "" {
		return style.ErrorStyle().Render("UNSET")
	}
	return value
}

func PrintTargetContext(context *state.TargetContext) {
	log.Info().
		Str("org", string(context.Org.Name)).
		Str("project", string(context.Project.Name)).
		Msg("Printing target context")
	writer := NewTabWriter()

	columns := []string{"ATTRIBUTE", "VALUE"}
	header := strings.Join(columns, "\t") + "\n"
	pcio.Fprint(writer, header)

	pcio.Fprintf(writer, "Org\t%s\n", labelUnsetIfEmpty(string(context.Org.Name)))
	pcio.Fprintf(writer, "Org ID\t%s\n", labelUnsetIfEmpty(string(context.Org.Id)))
	pcio.Fprintf(writer, "Project\t%s\n", labelUnsetIfEmpty(string(context.Project.Name)))
	pcio.Fprintf(writer, "Project ID\t%s\n", labelUnsetIfEmpty(string(context.Project.Id)))

	writer.Flush()
}
