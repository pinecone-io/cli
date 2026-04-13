package login

import (
	_ "embed"

	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/login"
	"github.com/spf13/cobra"
)

var (
	loginHelp = help.Long(`
		Authenticate with Pinecone via user login in a web browser.

		INTERACTIVE MODE (default)

		Opens a browser to the Pinecone login page and waits for you to complete
		authentication. The CLI automatically sets a default target organization
		and project. Use 'pc target' to change the target at any time.

		AGENTIC / NON-INTERACTIVE MODE (--json / -j, or non-TTY stdout)

		Uses a daemon-backed two-call flow designed for AI agents and scripts:

		First call — starts a background listener and returns immediately:
		  {"status":"pending","url":"<auth-url>","session_id":"<id>"}

		Open the URL in a browser to complete authentication. The background
		listener captures the OAuth callback automatically.

		Second call (or any other command) — completes the flow:
		  {"status":"authenticated","email":"...","org_id":"...","org_name":"...","project_id":"...","project_name":"..."}

		If the process is interrupted between calls, the background listener keeps
		running. The next invocation detects the pending session and resumes
		automatically. After authentication is complete, the first subsequent
		command also sets the target context automatically, so a separate
		'pc target' call is not required.
	`)
)

func NewLoginCmd() *cobra.Command {
	var jsonOutput bool
	var orgId string

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate with Pinecone via user login in a web browser",
		Long:  loginHelp,
		Example: help.Examples(`
			# Interactive login (opens a browser)
			pc login

			# Login scoped to a specific organization (enables SSO routing)
			pc login --org "ORG_ID"

			# Agentic login — first call returns a pending URL
			pc login --json

			# Agentic login — second call (or any command) completes the flow
			pc login --json
		`),
		GroupID: help.GROUP_AUTH.ID,
		Run: func(cmd *cobra.Command, args []string) {
			opts := login.Options{Json: jsonOutput}
			if cmd.Flags().Changed("org") {
				opts.OrgId = &orgId
			}
			login.Run(cmd.Context(), opts)
		},
	}

	cmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "emit JSON output")
	cmd.Flags().StringVar(&orgId, "org", "", "Organization ID to authenticate into (enables SSO routing for organizations with SSO enforced)")

	return cmd
}
