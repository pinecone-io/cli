package logout

import (
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/spf13/cobra"
)

func NewLogoutCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "logout",
		Short:   "Delete all saved tokens and keys",
		GroupID: help.GROUP_START.ID,
		Run: func(cmd *cobra.Command, args []string) {
			secrets.ConfigFile.Clear()
			fmt.Println(style.SuccessMsg("Secrets cleared."))
		},
	}

	return cmd
}
