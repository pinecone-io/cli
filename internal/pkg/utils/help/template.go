package help

import (
	"fmt"

	"github.com/pinecone-io/cli/internal/pkg/utils/docslinks"
	"github.com/pinecone-io/cli/internal/pkg/utils/style"
)

// Context: Cobra uses the text/template package to render help text.
// If you want to customize the help text, you have to tweak the default
// template. Unfortunately, it is very difficult to reason about.
// To do some styling things I wanted to interpolate some values, so
// I split it up into several smaller fragments just to make it easier
// to keep track of the order of positional arguments to the sprintf function.

var dot = style.Emphasis(style.Dot)

var description = `{{with (or .Long .Short)}}{{pcBlock .}}{{end}}

`

var usage = fmt.Sprintf(`%s{{if .Runnable}}
{{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
    {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}`, style.Heading("Usage"))

var aliasesAndExamples = fmt.Sprintf(`

%s
{{.NameAndAliases}}{{end}}{{if .HasExample}}

%s
{{pcExamples .Example}}{{end}}{{if .HasAvailableSubCommands}}{{$cmds := .Commands}}{{if eq (len .Groups) 0}}

`, style.Heading("Aliases"), style.Heading("Examples"))

var noGroupCmds = fmt.Sprintf(`Available Commands{{range $cmds}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  %s {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{else}}{{range $group := .Groups}}

`, dot)

var groupCmds = fmt.Sprintf(`{{.Title}}{{range $cmds}}{{if (and (eq .GroupID $group.ID) (or .IsAvailableCommand (eq .Name "help")))}}
  %s {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if not .AllChildCommandsHaveGroup}}

`, dot)

var additionalCmds = fmt.Sprintf(`%s{{range $cmds}}{{if (and (eq .GroupID "") (or .IsAvailableCommand (eq .Name "help")))}}
  %s {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

`, style.Heading("Additional Commands"), dot)

var flagsAndFooter = fmt.Sprintf(`%s
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

%s
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}

For in-depth documentation and resources, visit %s
`, style.Heading("Flags"), style.Heading("Global Flags"), style.URL(docslinks.DocsHome))

var HelpTemplate = description +
	usage +
	aliasesAndExamples +
	noGroupCmds +
	groupCmds +
	additionalCmds +
	flagsAndFooter
