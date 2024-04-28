package org

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/pinecone-io/cli/internal/pkg/dashboard"
	"github.com/spf13/cobra"

	"github.com/pinecone-io/cli/internal/pkg/utils/exit"
	"github.com/pinecone-io/cli/internal/pkg/utils/text"
)

type ListOrgCmdOptions struct {
	json bool
}

func NewListOrgsCmd() *cobra.Command {
	options := ListOrgCmdOptions{}

	cmd := &cobra.Command{
		Use:   "list <command>",
		Short: "list organizations",
		Run: func(cmd *cobra.Command, args []string) {
			orgs, err := dashboard.GetOrganizations()
			if err != nil {
				exit.Error(err)
			}

			if options.json {
				text.PrettyPrintJSON(orgs)
				return
			}

			printTable(orgs.Organizations)
		},
	}

	cmd.Flags().BoolVar(&options.json, "json", false, "output as JSON")

	return cmd
}

func printTable(orgs []dashboard.Organization) {
	writer := tabwriter.NewWriter(os.Stdout, 10, 1, 3, ' ', 0)

	columns := []string{"ID", "NAME", "PROJECTS"}
	header := strings.Join(columns, "\t") + "\n"
	fmt.Fprint(writer, header)

	for _, org := range orgs {
		values := []string{org.Id, org.Name, fmt.Sprintf("%d", len(org.Projects))}
		fmt.Fprintf(writer, strings.Join(values, "\t")+"\n")
	}
	writer.Flush()
}
