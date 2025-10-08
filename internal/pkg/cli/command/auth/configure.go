package auth

import (
	"context"
	"io"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/secrets"
	"github.com/pinecone-io/cli/internal/pkg/utils/configuration/state"
	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/help"
	"github.com/pinecone-io/cli/internal/pkg/utils/log"
	"github.com/pinecone-io/cli/internal/pkg/utils/msg"
	"github.com/pinecone-io/cli/internal/pkg/utils/oauth"
	"github.com/pinecone-io/cli/internal/pkg/utils/pcio"
	"github.com/pinecone-io/cli/internal/pkg/utils/presenters"
	"github.com/pinecone-io/cli/internal/pkg/utils/prompt"
	"github.com/pinecone-io/cli/internal/pkg/utils/sdk"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
	"github.com/pinecone-io/go-pinecone/v4/pinecone"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

type ConfigureCmdOptions struct {
	clientID            string
	clientSecret        string
	projectId           string
	apiKey              string
	readSecretFromStdin bool
	promptIfMissing     bool
	json                bool
}

var (
	configureHelp = help.Long(`
		Configure the CLI to use a service account or API key for authentication.
		
		When you configure a service account, the CLI automatically targets the organization
		associated with that account, and prompts you to select a project if multiple exist.
		
		An API overrides any explicitly targeted organization and project, instead targeting
		the organization and project associated with the API key itself. API keys do not grant
		Admin API access.
		
		See: http://docs.pinecone.io/reference/tools/cli-authentication
	`)
)

func NewConfigureCmd() *cobra.Command {
	options := ConfigureCmdOptions{}

	cmd := &cobra.Command{
		Use:   "configure",
		Short: "Configure authentication credentials for the Pinecone CLI",
		Long:  configureHelp,
		Example: help.Examples(`
			# Configure service account credentials
			pc auth configure --client-id "client-id" --client-secret "client-secret"

			# Configure default API key
			pc auth configure --api-key "api-key"
		`),
		GroupID: help.GROUP_AUTH.ID,
		Run: func(cmd *cobra.Command, args []string) {
			out := cmd.OutOrStdout()
			if quiet, _ := cmd.Flags().GetBool("quiet"); quiet {
				out = io.Discard
			}

			Run(cmd.Context(), IO{
				In:  cmd.InOrStdin(),
				Out: out,
				Err: cmd.ErrOrStderr(),
			}, options)
		},
	}

	cmd.Flags().StringVar(&options.clientID, "client-id", "", "Service account client id for the Pinecone CLI")
	cmd.Flags().StringVar(&options.clientSecret, "client-secret", "", "Service account client secret for the Pinecone CLI")
	cmd.Flags().StringVarP(&options.projectId, "project-id", "p", "", "The id of the project to target after authenticating with service account credentials")
	cmd.Flags().StringVar(&options.apiKey, "api-key", "", "Default API key override for the Pinecone CLI")
	cmd.Flags().BoolVar(&options.readSecretFromStdin, "client-secret-stdin", false, "Read the client secret from stdin")
	cmd.Flags().BoolVar(&options.promptIfMissing, "prompt-if-missing", false, "Prompt for missing credentials if not provided")
	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	return cmd
}

type IO struct {
	In  io.Reader
	Out io.Writer
	Err io.Writer
}

func Run(ctx context.Context, io IO, opts ConfigureCmdOptions) {
	clientID := strings.TrimSpace(opts.clientID)
	clientSecret := strings.TrimSpace(opts.clientSecret)
	globalAPIKey := strings.TrimSpace(opts.apiKey)

	// If clientSecret is not provided via options, prompt if needed
	if clientSecret == "" {
		if opts.readSecretFromStdin {
			secretBytes, err := ioReadAll(io.In)
			if err != nil {
				log.Error().Err(err).Msg("Error reading client secret from stdin")
				exit.Error(pcio.Errorf("error reading client secret from stdin: %w", err))
			}
			clientSecret = string(secretBytes)
		} else if opts.promptIfMissing && isTerminal(os.Stdin) {
			pcio.Fprint(io.Out, "Client Secret: ")
			secretBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				log.Error().Err(err).Msg("Error reading client secret from terminal")
				exit.Error(pcio.Errorf("error reading client secret from terminal: %w", err))
			}
			clientSecret = string(secretBytes)
		}
	}

	// If client_id is provided without a client_secret, error
	if clientID != "" && clientSecret == "" {
		log.Error().Msg("Error configuring authentication credentials")
		if !opts.json {
			msg.FailMsg("Client secret is required (use %s or %s to provide it)", style.Emphasis("--client-secret"), style.Emphasis("--client-secret-stdin"))
		}
		exit.Error(pcio.Errorf("client secret is required"))
		return
	}

	// If client_id and client_secret are configured, we need to use the AdminClient to fetch organization and project information for the service account
	if clientID != "" && clientSecret != "" {
		// Clear any existing user token login
		oauth.Logout()

		secrets.ClientId.Set(clientID)
		secrets.ClientSecret.Set(clientSecret)

		// Use Admin API to fetch organization and project information for the service account
		// so that we can set the target context, or allow the user to set it like they do through the login or target flow
		ac := sdk.NewPineconeAdminClient()

		// There should only be one organization listed for a service account
		orgs, err := ac.Organization.List(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Error listing service account organizations")
			exit.Error(pcio.Errorf("Error listing service account organizations: %w", err))
		}

		if len(orgs) == 0 {
			log.Error().Msg("No organizations found for service account")
			exit.ErrorMsg("No organizations found for service account")
		}

		targetOrg := orgs[0]

		state.TargetOrg.Set(state.TargetOrganization{
			Name: targetOrg.Name,
			Id:   targetOrg.Id,
		})
		if !opts.json {
			msg.SuccessMsg("Target organization set to %s", style.Emphasis(targetOrg.Name))
		}

		// List projects, and allow the user to pick one, or match the project-id if provided through the command
		projects, err := ac.Project.List(ctx)
		if err != nil {
			log.Error().Err(err).Msg("Error listing projects for service account")
			exit.Error(pcio.Errorf("Error listing projects for service account: %w", err))
		}

		var targetProject *pinecone.Project

		//  If the user has no projects, they can create one by running the project create command
		if len(projects) == 0 {
			log.Info().Msg("No projects found for service account")
			exit.SuccessMsg(pcio.Sprintf("No projects found for service account, you can create a project by running %s", style.Code("pc project create")))
		}

		// If the user has one project, set it as the target project
		if len(projects) == 1 {
			targetProject = projects[0]
			state.TargetProj.Set(state.TargetProject{
				Name: targetProject.Name,
				Id:   targetProject.Id,
			})
			if !opts.json {
				msg.SuccessMsg("Target project set to %s", style.Emphasis(targetProject.Name))
			}
			exit.Success()
		}

		// If there are multiple projects, select based on project-id, or allow the user to select one
		if opts.projectId != "" {
			for _, project := range projects {
				if project.Id == opts.projectId {
					targetProject = project
					break
				}
			}
		} else {
			targetProject = uiProjectSelector(projects)
		}

		state.TargetProj.Set(state.TargetProject{
			Name: targetProject.Name,
			Id:   targetProject.Id,
		})
		if !opts.json {
			msg.SuccessMsg("Target project set to %s", style.Emphasis(targetProject.Name))
		}

		// Update target credentials context
		state.TargetCreds.Set(state.TargetUser{
			AuthContext: state.AuthServiceAccount,
			Email:       "",
		})

		// Log out and clear oauth2 token if previously logged in and we've configured a service account
		oauth2Token := secrets.GetOAuth2Token()
		if oauth2Token.AccessToken != "" {
			oauth.Logout()
		}
	}

	// If a default API key has been configured, store it and update the target credentials context
	// This will override the AuthContext: state.AuthServiceAccount if set previously
	if globalAPIKey != "" {
		secrets.DefaultAPIKey.Set(globalAPIKey)
		state.TargetCreds.Set(state.TargetUser{
			AuthContext: state.AuthGlobalAPIKey,
			// Redact API key for presentational layer
			GlobalAPIKey: presenters.MaskHeadTail(globalAPIKey, 4, 4),
			Email:        "",
		})
	}

	// Output JSON if the option was passed
	if opts.json {
		json := text.IndentJSON(state.GetTargetContext())
		pcio.Println(json)
		return
	}

	pcio.Println()
	presenters.PrintTargetContext(state.GetTargetContext())
}

func ioReadAll(r io.Reader) ([]byte, error) {
	if r == nil {
		return []byte{}, nil
	}
	var buf strings.Builder
	tmp := make([]byte, 4096)
	for {
		n, err := r.Read(tmp)
		if n > 0 {
			buf.Write(tmp[:n])
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
	}
	return []byte(buf.String()), nil
}

func isTerminal(f *os.File) bool {
	if f == nil {
		return false
	}
	return term.IsTerminal(int(f.Fd()))
}

func uiProjectSelector(projects []*pinecone.Project) *pinecone.Project {
	var targetProject *pinecone.Project
	var targetProjectName string

	projectItems := []string{}
	projectMap := make(map[string]*pinecone.Project)
	for _, project := range projects {
		projectItems = append(projectItems, project.Name)
		projectMap[project.Name] = project
	}

	m2 := prompt.NewList(projectItems, len(projectItems)+6, "Choose a project to target", func() {
		pcio.Println("Exiting without targeting a project.")
		pcio.Printf("You can always run %s to set or change a project context later.\n", style.Code("pc target"))
		exit.Success()
	}, func(choice string) string {
		targetProjectName = choice
		return "Target project: " + choice
	})
	if _, err := tea.NewProgram(m2).Run(); err != nil {
		pcio.Println("Error running program:", err)
		os.Exit(1)
	}

	targetProject = projectMap[targetProjectName]
	return targetProject
}
