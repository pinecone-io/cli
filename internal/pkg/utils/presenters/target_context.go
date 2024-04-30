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
		return style.StatusRed("UNSET")
	}
	return value
}

func PrintTargetContext(context *state.TargetContext) {
	log.Info().Str("api", context.Api).Str("org", context.Org).Str("project", context.Project).Msg("Printing target context")
	writer := NewTabWriter()

	columns := []string{"ATTRIBUTE", "VALUE"}
	header := strings.Join(columns, "\t") + "\n"
	pcio.Fprint(writer, header)

	pcio.Fprintf(writer, "Api\t%s\n", labelUnsetIfEmpty(style.Emphasis(context.Api)))
	pcio.Fprintf(writer, "Org\t%s\n", labelUnsetIfEmpty(context.Org))
	pcio.Fprintf(writer, "Project\t%s\n", labelUnsetIfEmpty(context.Project))

	writer.Flush()
}
