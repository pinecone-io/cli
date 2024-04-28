package presenters

import (
	"fmt"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
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
	fmt.Fprint(writer, header)

	fmt.Fprintf(writer, "Api\t%s\n", labelUnsetIfEmpty(style.Emphasis(context.Api)))
	fmt.Fprintf(writer, "Org\t%s\n", labelUnsetIfEmpty(context.Org))
	fmt.Fprintf(writer, "Project\t%s\n", labelUnsetIfEmpty(context.Project))

	writer.Flush()
}
