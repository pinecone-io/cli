package presenters

import (
	"fmt"
	"strings"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
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
	writer := NewTabWriter()
	if context == nil {
		PrintEmptyState(writer, "target context")
		return
	}

	log.Info().
		Str("org", context.Organization.Name).
		Str("project", context.Project.Name).
		Msg("Printing target context")

	columns := []string{"ATTRIBUTE", "VALUE"}
	header := strings.Join(columns, "\t") + "\n"
	fmt.Fprint(writer, header)

	// Get API key for presentational layer
	defaultAPIKeyMasked := MaskHeadTail(secrets.DefaultAPIKey.Get(), 4, 4)

	fmt.Fprintf(writer, "Organization\t%s\n", labelUnsetIfEmpty(context.Organization.Name))
	fmt.Fprintf(writer, "Organization ID\t%s\n", labelUnsetIfEmpty(context.Organization.Id))
	fmt.Fprintf(writer, "Project\t%s\n", labelUnsetIfEmpty(context.Project.Name))
	fmt.Fprintf(writer, "Project ID\t%s\n", labelUnsetIfEmpty(context.Project.Id))
	fmt.Fprintf(writer, "Default API Key\t%s\n", labelUnsetIfEmpty(defaultAPIKeyMasked))

	writer.Flush()
}
