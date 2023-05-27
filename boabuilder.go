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

// WithProfiles is used to add any number of options to the boa Command
func (b *BoaCmdBuilder) WithProfiles(profs ...Profile) *BoaCmdBuilder {
	b.cmd.Profiles = append(b.cmd.Profiles, profs...)
	return b
}

// WithOptionsAndTemplate is used to add any number of options to the boa
// Command and applies the WithOptionsTemplate() method as well.
func (b *BoaCmdBuilder) WithOptionsAndTemplate(opts ...Option) *BoaCmdBuilder {
	return b.WithOptions(opts...).WithOptionsTemplate()
}

// WithUsageTemplate is used to add a custom template for usage text
func (b *BoaCmdBuilder) WithUsageTemplate(template string) *BoaCmdBuilder {
	cmd := b.Build()
	b.WithUsageFunc(cmd.UsageFunc(template))
	return b
}

// WithHelpTemplate is used to add a custom template for help text
func (b *BoaCmdBuilder) WithHelpTemplate(template string) *BoaCmdBuilder {
	cmd := b.Build()
	b.WithHelpFunc(cmd.HelpFunc(template))
	return b
}

// WithOptionsTemplate is used to add options to the usage and help text
func (b *BoaCmdBuilder) WithOptionsTemplate() *BoaCmdBuilder {
	template := b.Build().OptionsTemplate()
	return b.WithUsageTemplate(template).WithHelpTemplate(template)
}

// WithValidArgsFromOptions updates the underlying cobra.Command's ValidArgs
// with all arguments from the boa Commands options
//
// Should be used in conjunction with a CobraCmdBuilder
//
//	WithArgs(cobra.OnlyValidArgs)
func (b *BoaCmdBuilder) WithValidArgsFromOptions() *BoaCmdBuilder {
	for _, opt := range b.cmd.Opts {
		b.cmd.ValidArgs = append(b.cmd.ValidArgs, opt.Args...)
	}
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
