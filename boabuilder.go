package boa

import (
	"github.com/spf13/cobra"
)

// BoaCmdBuilder is a wrapper for the CobraCmdBuilder that allows for building
// boa Commands using the same builder methods as a CobraCmdBuilder, but with
// additional builder methods specific to a boa Command.
type BoaCmdBuilder struct {
	*CobraCmdBuilder
	cmd *Command
}

// ToBoaCmdBuilder is used to convert a cobra.Command to a BoaCmdBuilder.
func ToBoaCmdBuilder(cmd *cobra.Command) *BoaCmdBuilder {
	return &BoaCmdBuilder{
		&CobraCmdBuilder{cmd},
		&Command{cmd, []Option{}, []Profile{}},
	}
}

// NewCmd creates a new BoaCmdBuilder and sets the use for the underlying cobra
// Command.
func NewCmd(use string) *BoaCmdBuilder {
	cobraBuilder := NewCobraCmd(use)
	return &BoaCmdBuilder{
		CobraCmdBuilder: cobraBuilder,
		cmd: &Command{
			Command: cobraBuilder.Build(),
			Opts:    []Option{},
		},
	}
}

// WithOptions is used to add any number of options to the boa Command
func (b *BoaCmdBuilder) WithOptions(opts ...Option) *BoaCmdBuilder {
	b.cmd.Opts = append(b.cmd.Opts, opts...)
	return b
}

// WithValidOptions is used to add any number of options to the boa Command and
// set them as ValidArgs
func (b *BoaCmdBuilder) WithValidOptions(opts ...Option) *BoaCmdBuilder {
	b.cmd.Opts = append(b.cmd.Opts, opts...)
	for _, opt := range b.cmd.Opts {
		b.cmd.ValidArgs = append(b.cmd.ValidArgs, opt.Args...)
	}
	return b
}

// WithProfiles is used to add any number of options to the boa Command
func (b *BoaCmdBuilder) WithProfiles(profs ...Profile) *BoaCmdBuilder {
	b.cmd.Profiles = append(b.cmd.Profiles, profs...)
	return b
}

// WithValidProfiles is used to add any number of options to the boa Command and
// set them as ValidArgs
func (b *BoaCmdBuilder) WithValidProfiles(profs ...Profile) *BoaCmdBuilder {
	b.cmd.Profiles = append(b.cmd.Profiles, profs...)
	for _, prof := range b.cmd.Profiles {
		b.cmd.ValidArgs = append(b.cmd.ValidArgs, prof.Args...)
	}
	return b
}

// WithUsageTemplate is used to add a custom template for usage text
func (b *BoaCmdBuilder) WithUsageTemplate(template string) *BoaCmdBuilder {
	b.WithUsageFunc(b.cmd.UsageFunc(template))
	return b
}

// WithHelpTemplate is used to add a custom template for help text
func (b *BoaCmdBuilder) WithHelpTemplate(template string) *BoaCmdBuilder {
	b.WithHelpFunc(b.cmd.HelpFunc(template))
	return b
}

// WithOptionsTemplate is used to add options to the usage and help text
func (b *BoaCmdBuilder) WithOptionsTemplate() *BoaCmdBuilder {
	template := b.cmd.OptionsTemplate()
	return b.WithUsageTemplate(template).WithHelpTemplate(template)
}

// WithMinValidArgs will cause the command to throw an error if at least minArgs
// valid arguments are not provided
func (b *BoaCmdBuilder) WithMinValidArgs(minArgs int) *BoaCmdBuilder {
	b.cmd.Args = cobra.MatchAll(cobra.MinimumNArgs(minArgs), cobra.OnlyValidArgs)
	return b
}

// WithMaxValidArgs will cause the command to throw an error if more than
// maxArgs valid arguments are provided
func (b *BoaCmdBuilder) WithMaxValidArgs(maxArgs int) *BoaCmdBuilder {
	b.cmd.Args = cobra.MatchAll(cobra.MaximumNArgs(maxArgs), cobra.OnlyValidArgs)
	return b
}

// ToCobraCmdBuilder returns a CobraCmdBuilder from a BoaCmdBuilder
//
// This method isn't particularly useful as a BoaCmdBuilder is also a
// CobraCmdBuilder and has access to all CobraCmdBuilder methods; however, it
// does give the user the choice to make it more apparent to others that
// subsequent methods are specific to a CobraCmdBuilder.
func (b *BoaCmdBuilder) ToCobraCmdBuilder() *CobraCmdBuilder {
	return b.CobraCmdBuilder
}

// BuildCobraCmd returns a cobra.Command from a BoaCmdBuilder
//
// This method allows bypassing the ToCobraCmdBuilder() step before Build()
func (b *BoaCmdBuilder) BuildCobraCmd() *cobra.Command {
	return b.cmd.Command
}

// Build returns a boa Command from a BoaCmdBuilder
func (b *BoaCmdBuilder) Build() *Command {
	return b.cmd
}
