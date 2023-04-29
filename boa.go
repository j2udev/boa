package boa

import (
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

type (
	// Option is used to define multiple positional args in which the positional
	// args can have a description. Aliases for the args can be added to the Args
	// slice.
	Option struct {
		Args []string
		Desc string
	}
	// Command is a wrapper for the cobra Command that adds additional fields to
	// support better usage, help, etc.
	Command struct {
		*cobra.Command
		Opts []Option
	}
)

// Build returns a boa Command from a BoaCmdBuilder
func (b Command) ToBuilder() *BoaCmdBuilder {
	return &BoaCmdBuilder{
		NewCobraCmd(b.Use),
		&b,
	}
}

// UsageFunc overrides the default UsageFunc used by boa to facilitate showing
// a custom usage template
func (c Command) UsageFunc(template string) func(*cobra.Command) error {
	return func(cmd *cobra.Command) error {
		w := tabwriter.NewWriter(os.Stdout, 8, 8, 8, ' ', 0)
		err := tmpl(w, template, c)
		if err != nil {
			cmd.PrintErrln(err)
		}
		return err
	}
}

// HelpFunc overrides the default HelpFunc used by cobra to facilitate showing
// a custom help template
func (c Command) HelpFunc(template string) func(*cobra.Command, []string) {
	return func(cmd *cobra.Command, s []string) {
		w := tabwriter.NewWriter(os.Stdout, 3, 3, 3, ' ', 0)
		err := tmpl(w, template, c)
		if err != nil {
			cmd.PrintErrln(err)
		}
	}
}

// OptionsTemplate is used to override the cobra UsageTemplate to facilitate
// options and other CLI parameters
func (c Command) OptionsTemplate() string {
	return `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasOptions}} [options]{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}{{$cmds := .Commands}}{{if eq (len .Groups) 0}}

Available Commands:{{range $cmds}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{else}}{{range $group := .Groups}}

{{.Title}}{{range $cmds}}{{if (and (eq .GroupID $group.ID) (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if not .AllChildCommandsHaveGroup}}

Additional Commands:{{range $cmds}}{{if (and (eq .GroupID "") (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{end}}{{end}}{{if .HasOptions}}

Options:{{range .Opts }}
  {{.Args | sliceToCsv}}	{{.Desc}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`
}

// HasOptions returns whether the boa Command has any options defined; this is
// primary used for templating purposes.
func (c Command) HasOptions() bool {
	if c.Opts == nil || len(c.Opts) == 0 {
		return false
	}
	return true
}
