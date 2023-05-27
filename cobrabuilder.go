package boa

import (
	"net"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// CobraCmdBuilder is a builder for cobra.Command fields and chaining other
// helpful methods. Flags can be added to a command using builder methods as
// well.
type CobraCmdBuilder struct {
	cmd *cobra.Command
}

// ToCobraCmdBuilder is used to convert an existing cobra.Command to a
// CobraCmdBuilder.
func ToCobraCmdBuilder(cmd *cobra.Command) *CobraCmdBuilder {
	return &CobraCmdBuilder{cmd}
}

// NewCobraCmd creates a new CobraCmdBuilder and sets the use for the
// underlying cobra.Command
//
// Use is the one-line usage message.
// Recommended syntax is as follows:
//
//	[ ] identifies an optional argument. Arguments that are not enclosed in brackets are required.
//	... indicates that you can specify multiple values for the previous argument.
//	|   indicates mutually exclusive information. You can use the argument to the left of the separator or the
//	    argument to the right of the separator. You cannot use both arguments in a single use of the command.
//	{ } delimits a set of mutually exclusive arguments when one of the arguments is required. If the arguments are
//	    optional, they are enclosed in brackets ([ ]).
//
// Example: add [-F file | -D dir]... [-f format] profile
func NewCobraCmd(use string) *CobraCmdBuilder {
	return &CobraCmdBuilder{
		cmd: &cobra.Command{
			Use: use,
		},
	}
}

// WithAliases is an array of aliases that can be used instead of the first word
// in Use.
func (b *CobraCmdBuilder) WithAliases(aliases []string) *CobraCmdBuilder {
	b.cmd.Aliases = aliases
	return b
}

// SuggestFor is an array of command names for which this command will be
// suggested - similar to aliases but only suggests.
func (b *CobraCmdBuilder) SuggestFor(cmds []string) *CobraCmdBuilder {
	b.cmd.SuggestFor = cmds
	return b
}

// WithShortDescription is the short description shown in the 'help' output.
func (b *CobraCmdBuilder) WithShortDescription(short string) *CobraCmdBuilder {
	b.cmd.Short = short
	return b
}

// WithGroupId is the group id under which this subcommand is grouped in the
// 'help' output of its parent.
func (b *CobraCmdBuilder) WithGroupID(groupId string) *CobraCmdBuilder {
	b.cmd.GroupID = groupId
	return b
}

// WithLongDescription is the long message shown in the 'help <this-command>'
// output.
func (b *CobraCmdBuilder) WithLongDescription(long string) *CobraCmdBuilder {
	b.cmd.Long = long
	return b
}

// WithExample is examples of how to use the command.
func (b *CobraCmdBuilder) WithExample(example string) *CobraCmdBuilder {
	b.cmd.Example = example
	return b
}

// WithValidArgs is list of all valid non-flag arguments that are accepted in
// shell completions
func (b *CobraCmdBuilder) WithValidArgs(validArgs []string) *CobraCmdBuilder {
	b.cmd.ValidArgs = append(b.cmd.ValidArgs, validArgs...)
	return b
}

// WithValidArgsFunction is an optional function that provides valid non-flag
// arguments for shell completion. It is a dynamic version of using ValidArgs.
// Only one of ValidArgs and ValidArgsFunction can be used for a command.
func (b *CobraCmdBuilder) WithValidArgsFunction(validArgsFunc func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective)) *CobraCmdBuilder {
	b.cmd.ValidArgsFunction = validArgsFunc
	return b
}

// WithArgs sets the expected arguments for the command.
//
// For example:
//
//	WithArgs(cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs))
func (b *CobraCmdBuilder) WithArgs(args cobra.PositionalArgs) *CobraCmdBuilder {
	b.cmd.Args = args
	return b
}

// WithArgAliases is List of aliases for ValidArgs. These are not suggested to
// the user in the shell completion, but accepted if entered manually.
func (b *CobraCmdBuilder) WithArgAliases(argAliases []string) *CobraCmdBuilder {
	b.cmd.ArgAliases = argAliases
	return b
}

// WithBashCompletionFunction is custom bash functions used by the legacy bash
// autocompletion generator. For portability with other shells, it is
// recommended to instead use ValidArgsFunction
func (b *CobraCmdBuilder) WithBashCompletionFunction(bashCompletionFunction string) *CobraCmdBuilder {
	b.cmd.BashCompletionFunction = bashCompletionFunction
	return b
}

// Deprecated defines if this command is deprecated and should print this string
// when used
func (b *CobraCmdBuilder) Deprecated(deprecated string) *CobraCmdBuilder {
	b.cmd.Deprecated = deprecated
	return b
}

// WithAnnotations are key/value pairs that can be used by applications to
// identify or group commands.
func (b *CobraCmdBuilder) WithAnnotations(annotations map[string]string) *CobraCmdBuilder {
	b.cmd.Annotations = annotations
	return b
}

// Version defines the version for this command. If this value is non-empty and
// the command does not define a "version" flag, a "version" boolean flag will
// be added to the command and, if specified, will print content of the
// "Version" variable. A shorthand "v" flag will also be added if the command
// does not define one.
func (b *CobraCmdBuilder) WithVersion(version string) *CobraCmdBuilder {
	b.cmd.Version = version
	return b
}

func (b *CobraCmdBuilder) WithNoOp() *CobraCmdBuilder {
	return b.WithRunFunc(func(*cobra.Command, []string) {})
}

// The *Run functions are executed in the following order:
//   - PersistentPreRun()
//   - PreRun()
//   - Run()
//   - PostRun()
//   - PersistentPostRun()
//
// All functions get the same args, the arguments after the command name.
//
// WithPersistentPreRunFunc: children of this command will inherit and execute.
func (b *CobraCmdBuilder) WithPersistentPreRunFunc(f func(cmd *cobra.Command, args []string)) *CobraCmdBuilder {
	b.cmd.PersistentPreRun = f
	return b
}

// The *Run functions are executed in the following order:
//   - PersistentPreRun()
//   - PreRun()
//   - Run()
//   - PostRun()
//   - PersistentPostRun()
//
// All functions get the same args, the arguments after the command name.
//
// WithPersistentPreRunEFunc: PersistentPreRun but returns an error.
func (b *CobraCmdBuilder) WithPersistentPreRunEFunc(f func(cmd *cobra.Command, args []string) error) *CobraCmdBuilder {
	b.cmd.PersistentPreRunE = f
	return b
}

// The *Run functions are executed in the following order:
//   - PersistentPreRun()
//   - PreRun()
//   - Run()
//   - PostRun()
//   - PersistentPostRun()
//
// All functions get the same args, the arguments after the command name.
//
// WithPreRunFunc: children of this command will not inherit.
func (b *CobraCmdBuilder) WithPreRunFunc(f func(cmd *cobra.Command, args []string)) *CobraCmdBuilder {
	b.cmd.PreRun = f
	return b
}

// The *Run functions are executed in the following order:
//   - PersistentPreRun()
//   - PreRun()
//   - Run()
//   - PostRun()
//   - PersistentPostRun()
//
// All functions get the same args, the arguments after the command name.
//
// WithPreRunEFunc: PreRun but returns an error.
func (b *CobraCmdBuilder) WithPreRunEFunc(f func(cmd *cobra.Command, args []string) error) *CobraCmdBuilder {
	b.cmd.PreRunE = f
	return b
}

// The *Run functions are executed in the following order:
//   - PersistentPreRun()
//   - PreRun()
//   - Run()
//   - PostRun()
//   - PersistentPostRun()
//
// All functions get the same args, the arguments after the command name.
//
// WithRunFunc: Typically the actual work function. Most commands will only
// implement this.
func (b *CobraCmdBuilder) WithRunFunc(f func(cmd *cobra.Command, args []string)) *CobraCmdBuilder {
	b.cmd.Run = f
	return b
}

// The *Run functions are executed in the following order:
//   - PersistentPreRun()
//   - PreRun()
//   - Run()
//   - PostRun()
//   - PersistentPostRun()
//
// All functions get the same args, the arguments after the command name.
//
// WithRunEFunc: Run but returns an error.
func (b *CobraCmdBuilder) WithRunEFunc(f func(cmd *cobra.Command, args []string) error) *CobraCmdBuilder {
	b.cmd.RunE = f
	return b
}

// The *Run functions are executed in the following order:
//   - PersistentPreRun()
//   - PreRun()
//   - Run()
//   - PostRun()
//   - PersistentPostRun()
//
// All functions get the same args, the arguments after the command name.
//
// WithPostRunFunc: run after the Run command.
func (b *CobraCmdBuilder) WithPostRunFunc(f func(cmd *cobra.Command, args []string)) *CobraCmdBuilder {
	b.cmd.PostRun = f
	return b
}

// The *Run functions are executed in the following order:
//   - PersistentPreRun()
//   - PreRun()
//   - Run()
//   - PostRun()
//   - PersistentPostRun()
//
// All functions get the same args, the arguments after the command name.
//
// WithPostRunEFunc: PostRun but returns an error.
func (b *CobraCmdBuilder) WithPostRunEFunc(f func(cmd *cobra.Command, args []string) error) *CobraCmdBuilder {
	b.cmd.PostRunE = f
	return b
}

// The *Run functions are executed in the following order:
//   - PersistentPreRun()
//   - PreRun()
//   - Run()
//   - PostRun()
//   - PersistentPostRun()
//
// All functions get the same args, the arguments after the command name.
//
// WithPersistentPostRunFunc: children of this command will inherit and execute
// after PostRun.
func (b *CobraCmdBuilder) WithPersistentPostRunFunc(f func(cmd *cobra.Command, args []string)) *CobraCmdBuilder {
	b.cmd.PersistentPostRun = f
	return b
}

// The *Run functions are executed in the following order:
//   - PersistentPreRun()
//   - PreRun()
//   - Run()
//   - PostRun()
//   - PersistentPostRun()
//
// All functions get the same args, the arguments after the command name.
//
// WithPersistentPostRunEFunc: PersistentPostRun but returns an error.
func (b *CobraCmdBuilder) WithPersistentPostRunEFunc(f func(cmd *cobra.Command, args []string) error) *CobraCmdBuilder {
	b.cmd.PersistentPostRunE = f
	return b
}

// WithFParseErrWhitelist flag parse errors to be ignored.
func (b *CobraCmdBuilder) WithFParseErrWhitelist(flagParseErrors cobra.FParseErrWhitelist) *CobraCmdBuilder {
	b.cmd.FParseErrWhitelist = flagParseErrors
	return b
}

// WithCompletionOptions is a set of options to control the handling of shell
// completion.
func (b *CobraCmdBuilder) WithCompletionOptions(options cobra.CompletionOptions) *CobraCmdBuilder {
	b.cmd.CompletionOptions = options
	return b
}

// TraverseChildren parses flags on all parents before executing child command.
func (b *CobraCmdBuilder) TraverseChildren() *CobraCmdBuilder {
	b.cmd.TraverseChildren = true
	return b
}

// Hidden defines if this command is hidden and should NOT show up in the list
// of available commands.
func (b *CobraCmdBuilder) Hidden() *CobraCmdBuilder {
	b.cmd.Hidden = true
	return b
}

// SilenceErrors is an option to quiet errors down stream.
func (b *CobraCmdBuilder) SilenceErrors() *CobraCmdBuilder {
	b.cmd.SilenceErrors = true
	return b
}

// SilenceUsage is an option to silence usage when an error occurs.
func (b *CobraCmdBuilder) SilenceUsage() *CobraCmdBuilder {
	b.cmd.SilenceUsage = true
	return b
}

// DisableFlagParsing DisableFlagParsing disables the flag parsing. If this is
// true all flags will be passed to the command as arguments.
func (b *CobraCmdBuilder) DisableFlagParsing() *CobraCmdBuilder {
	b.cmd.DisableFlagParsing = true
	return b
}

// DisableAutoGenTag defines, if gen tag ("Auto generated by spf13/cobra...")
// will be printed by generating docs for this command.
func (b *CobraCmdBuilder) DisableAutoGenTag() *CobraCmdBuilder {
	b.cmd.DisableAutoGenTag = true
	return b
}

// DisableFlagsInUseLine will disable the addition of [flags] to the usage line
// of a command when printing help or generating docs
func (b *CobraCmdBuilder) DisableFlagsInUseLine() *CobraCmdBuilder {
	b.cmd.DisableFlagsInUseLine = true
	return b
}

// DisableSuggestions disables the suggestions based on Levenshtein distance
// that go along with 'unknown command' messages.
func (b *CobraCmdBuilder) DisableSuggestions() *CobraCmdBuilder {
	b.cmd.DisableSuggestions = true
	return b
}

// WithSuggestionsMinimumDistance defines minimum levenshtein distance to
// display suggestions. Must be > 0.
func (b *CobraCmdBuilder) WithSuggestionsMinimumDistance(distance int) *CobraCmdBuilder {
	b.cmd.SuggestionsMinimumDistance = distance
	return b
}

// WithSubCommands adds one or more commands to this parent command.
func (b *CobraCmdBuilder) WithSubCommands(cmds ...*cobra.Command) *CobraCmdBuilder {
	b.cmd.AddCommand(cmds...)
	return b
}

// WithUsageTemplate sets usage template. Can be defined by Application.
func (b *CobraCmdBuilder) WithUsageTemplate(template string) *CobraCmdBuilder {
	b.cmd.SetUsageTemplate(template)
	return b
}

// WithHelpTemplate sets help template to be used. Application can use it to set
// custom template.
func (b *CobraCmdBuilder) WithHelpTemplate(template string) *CobraCmdBuilder {
	b.cmd.SetHelpTemplate(template)
	return b
}

// WithUsageFunc sets usage function. Usage can be defined by application.
func (b *CobraCmdBuilder) WithUsageFunc(function func(*cobra.Command) error) *CobraCmdBuilder {
	b.cmd.SetUsageFunc(function)
	return b
}

// WithHelpFunc sets help function. Can be defined by Application.
func (b *CobraCmdBuilder) WithHelpFunc(function func(*cobra.Command, []string)) *CobraCmdBuilder {
	b.cmd.SetHelpFunc(function)
	return b
}

// WithBoolFlag defines a bool flag with specified name, default value, and
// usage string. The return value is the address of a bool variable that stores
// the value of the flag.
func (b *CobraCmdBuilder) WithBoolFlag(name string, value bool, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Bool(name, value, usage)
	return b
}

// WithBoolPFlag BoolP is like Bool, but accepts a shorthand letter that can be
// used after a single dash.
func (b *CobraCmdBuilder) WithBoolPFlag(name string, shorthand string, value bool, usage string) *CobraCmdBuilder {
	b.cmd.Flags().BoolP(name, shorthand, value, usage)
	return b
}

// WithBoolVarFlag defines a bool flag with specified name, default value, and
// usage string. The argument p points to a bool variable in which to store the
// value of the flag.
func (b *CobraCmdBuilder) WithBoolVarFlag(variable *bool, name string, value bool, usage string) *CobraCmdBuilder {
	b.cmd.Flags().BoolVar(variable, name, value, usage)
	return b
}

// WithBoolVarPFlag is like BoolVar, but accepts a shorthand letter that can be
// used after a single dash.
func (b *CobraCmdBuilder) WithBoolVarPFlag(variable *bool, name string, shorthand string, value bool, usage string) *CobraCmdBuilder {
	b.cmd.Flags().BoolVarP(variable, name, shorthand, value, usage)
	return b
}

// WithBoolSliceFlag defines a []bool flag with specified name, default value,
// and usage string. The return value is the address of a []bool variable that
// stores the value of the flag.
func (b *CobraCmdBuilder) WithBoolSliceFlag(name string, value []bool, usage string) *CobraCmdBuilder {
	b.cmd.Flags().BoolSlice(name, value, usage)
	return b
}

// WithBoolSlicePFlag is like BoolSlice, but accepts a shorthand letter that
// can be used after a single dash.
func (b *CobraCmdBuilder) WithBoolSlicePFlag(name string, shorthand string, value []bool, usage string) *CobraCmdBuilder {
	b.cmd.Flags().BoolSliceP(name, shorthand, value, usage)
	return b
}

// WithBoolSliceVarFlag defines a boolSlice flag with specified name, default
// value, and usage string. The argument p points to a []bool variable in which
// to store the value of the flag.
func (b *CobraCmdBuilder) WithBoolSliceVarFlag(variable *[]bool, name string, value []bool, usage string) *CobraCmdBuilder {
	b.cmd.Flags().BoolSliceVar(variable, name, value, usage)
	return b
}

// WithBoolSliceVarPFlag is like BoolSliceVar, but accepts a shorthand letter
// that can be used after a single dash.
func (b *CobraCmdBuilder) WithBoolSliceVarPFlag(variable *[]bool, name string, shorthand string, value []bool, usage string) *CobraCmdBuilder {
	b.cmd.Flags().BoolSliceVarP(variable, name, shorthand, value, usage)
	return b
}

// WithBytesBase64Flag defines an []byte flag with specified name, default
// value, and usage string. The return value is the address of an []byte
// variable that stores the value of the flag.
func (b *CobraCmdBuilder) WithBytesBase64Flag(name string, value []byte, usage string) *CobraCmdBuilder {
	b.cmd.Flags().BytesBase64(name, value, usage)
	return b
}

// WithBytesBase64PFlag is like BytesBase64, but accepts a shorthand letter that
// can be used after a single dash.
func (b *CobraCmdBuilder) WithBytesBase64PFlag(name string, shorthand string, value []byte, usage string) *CobraCmdBuilder {
	b.cmd.Flags().BytesBase64P(name, shorthand, value, usage)
	return b
}

// WithBytesBase64VarFlag defines an []byte flag with specified name, default
// value, and usage string. The argument p points to an []byte variable in which
// to store the value of the flag.
func (b *CobraCmdBuilder) WithBytesBase64VarFlag(variable *[]byte, name string, value []byte, usage string) *CobraCmdBuilder {
	b.cmd.Flags().BytesBase64Var(variable, name, value, usage)
	return b
}

// WithBytesBase64VarPFlag BytesBase64VarP is like BytesBase64Var, but accepts a
// shorthand letter that can be used after a single dash.
func (b *CobraCmdBuilder) WithBytesBase64VarPFlag(variable *[]byte, name string, shorthand string, value []byte, usage string) *CobraCmdBuilder {
	b.cmd.Flags().BytesBase64VarP(variable, name, shorthand, value, usage)
	return b
}

// WithBytesHexFlag defines an []byte flag with specified name, default value,
// and usage string. The return value is the address of an []byte variable that
// stores the value of the flag.
func (b *CobraCmdBuilder) WithBytesHexFlag(name string, value []byte, usage string) *CobraCmdBuilder {
	b.cmd.Flags().BytesHex(name, value, usage)
	return b
}

// WithBytesHexPFlag is like BytesHex, but accepts a shorthand letter that can
// be used after a single dash.
func (b *CobraCmdBuilder) WithBytesHexPFlag(name string, shorthand string, value []byte, usage string) *CobraCmdBuilder {
	b.cmd.Flags().BytesHexP(name, shorthand, value, usage)
	return b
}

// WithBytesHexVarFlag BytesHexVar defines an []byte flag with specified name,
// default value, and usage string. The argument p points to an []byte variable
// in which to store the value of the flag.
func (b *CobraCmdBuilder) WithBytesHexVarFlag(variable *[]byte, name string, value []byte, usage string) *CobraCmdBuilder {
	b.cmd.Flags().BytesHexVar(variable, name, value, usage)
	return b
}

// WithBytesHexVarPFlag is like BytesHexVar, but accepts a shorthand letter that
// can be used after a single dash.
func (b *CobraCmdBuilder) WithBytesHexVarPFlag(variable *[]byte, name string, shorthand string, value []byte, usage string) *CobraCmdBuilder {
	b.cmd.Flags().BytesHexVarP(variable, name, shorthand, value, usage)
	return b
}

// WithCountFlag defines a count flag with specified name, default value, and
// usage string. The return value is the address of an int variable that stores
// the value of the flag. A count flag will add 1 to its value every time it is
// found on the command line.
func (b *CobraCmdBuilder) WithCountFlag(name string, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Count(name, usage)
	return b
}

// WithCountPFlag is like Count only takes a shorthand for the flag name.
func (b *CobraCmdBuilder) WithCountPFlag(name string, shorthand string, usage string) *CobraCmdBuilder {
	b.cmd.Flags().CountP(name, shorthand, usage)
	return b
}

// WithCountVarFlag defines a count flag with specified name, default value, and
// usage string. The argument p points to an int variable in which to store the
// value of the flag. A count flag will add 1 to its value every time it is
// found on the command line
func (b *CobraCmdBuilder) WithCountVarFlag(variable *int, name string, usage string) *CobraCmdBuilder {
	b.cmd.Flags().CountVar(variable, name, usage)
	return b
}

// WithCountVarPFlag is like CountVar only take a shorthand for the flag name.
func (b *CobraCmdBuilder) WithCountVarPFlag(variable *int, name string, shorthand string, usage string) *CobraCmdBuilder {
	b.cmd.Flags().CountVarP(variable, name, shorthand, usage)
	return b
}

// WithDurationFlag defines a time.Duration flag with specified name, default
// value, and usage string. The return value is the address of a time.Duration
// variable that stores the value of the flag.
func (b *CobraCmdBuilder) WithDurationFlag(name string, value time.Duration, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Duration(name, value, usage)
	return b
}

// WithDurationPFlag is like Duration, but accepts a shorthand letter that can
// be used after a single dash.
func (b *CobraCmdBuilder) WithDurationPFlag(name string, shorthand string, value time.Duration, usage string) *CobraCmdBuilder {
	b.cmd.Flags().DurationP(name, shorthand, value, usage)
	return b
}

// WithDurationVarFlag defines a time.Duration flag with specified name, default
// value, and usage string. The argument p points to a time.Duration variable in
// which to store the value of the flag.
func (b *CobraCmdBuilder) WithDurationVarFlag(variable *time.Duration, name string, value time.Duration, usage string) *CobraCmdBuilder {
	b.cmd.Flags().DurationVar(variable, name, value, usage)
	return b
}

// WithDurationVarPFlag is like DurationVar, but accepts a shorthand letter that
// can be used after a single dash.
func (b *CobraCmdBuilder) WithDurationVarPFlag(variable *time.Duration, name string, shorthand string, value time.Duration, usage string) *CobraCmdBuilder {
	b.cmd.Flags().DurationVarP(variable, name, shorthand, value, usage)
	return b
}

// WithDurationSliceFlag defines a []time.Duration flag with specified name,
// default value, and usage string. The return value is the address of a
// []time.Duration variable that stores the value of the flag.
func (b *CobraCmdBuilder) WithDurationSliceFlag(name string, value []time.Duration, usage string) *CobraCmdBuilder {
	b.cmd.Flags().DurationSlice(name, value, usage)
	return b
}

// WithDurationSlicePFlag is like DurationSlice, but accepts a shorthand letter
// that can be used after a single dash.
func (b *CobraCmdBuilder) WithDurationSlicePFlag(name string, shorthand string, value []time.Duration, usage string) *CobraCmdBuilder {
	b.cmd.Flags().DurationSliceP(name, shorthand, value, usage)
	return b
}

// WithDurationSliceVarFlag defines a durationSlice flag with specified name,
// default value, and usage string. The argument p points to a []time.Duration
// variable in which to store the value of the flag.
func (b *CobraCmdBuilder) WithDurationSliceVarFlag(variable *[]time.Duration, name string, value []time.Duration, usage string) *CobraCmdBuilder {
	b.cmd.Flags().DurationSliceVar(variable, name, value, usage)
	return b
}

// WithDurationSliceVarPFlag is like DurationSliceVar, but accepts a shorthand
// letter that can be used after a single dash.
func (b *CobraCmdBuilder) WithDurationSliceVarPFlag(variable *[]time.Duration, name string, shorthand string, value []time.Duration, usage string) *CobraCmdBuilder {
	b.cmd.Flags().DurationSliceVarP(variable, name, shorthand, value, usage)
	return b
}

// WithFloat32Flag defines a float32 flag with specified name, default value,
// and usage string. The return value is the address of a float32 variable that
// stores the value of the flag.
func (b *CobraCmdBuilder) WithFloat32Flag(name string, value float32, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Float32(name, value, usage)
	return b
}

// WithFloat32PFlag is like Float32, but accepts a shorthand letter that can be
// used after a single dash.
func (b *CobraCmdBuilder) WithFloat32PFlag(name string, shorthand string, value float32, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Float32P(name, shorthand, value, usage)
	return b
}

// WithFloat32VarFlag defines a float32 flag with specified name, default value,
// and usage string. The argument p points to a float32 variable in which to
// store the value of the flag.
func (b *CobraCmdBuilder) WithFloat32VarFlag(variable *float32, name string, value float32, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Float32Var(variable, name, value, usage)
	return b
}

// WithFloat32VarPFlag is like Float32Var, but accepts a shorthand letter that
// can be used after a single dash.
func (b *CobraCmdBuilder) WithFloat32VarPFlag(variable *float32, name string, shorthand string, value float32, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Float32VarP(variable, name, shorthand, value, usage)
	return b
}

// WithFloat32SliceFlag defines a []float32 flag with specified name, default
// value, and usage string. The return value is the address of a []float32
// variable that stores the value of the flag.
func (b *CobraCmdBuilder) WithFloat32SliceFlag(name string, value []float32, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Float32Slice(name, value, usage)
	return b
}

// WithFloat32SlicePFlag is like Float32Slice, but accepts a shorthand letter
// that can be used after a single dash.
func (b *CobraCmdBuilder) WithFloat32SlicePFlag(name string, shorthand string, value []float32, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Float32SliceP(name, shorthand, value, usage)
	return b
}

// WithFloat32SliceVarFlag defines a float32Slice flag with specified name,
// default value, and usage string. The argument p points to a []float32
// variable in which to store the value of the flag.
func (b *CobraCmdBuilder) WithFloat32SliceVarFlag(variable *[]float32, name string, value []float32, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Float32SliceVar(variable, name, value, usage)
	return b
}

// WithFloat32SliceVarPFlag is like Float32SliceVar, but accepts a shorthand
// letter that can be used after a single dash.
func (b *CobraCmdBuilder) WithFloat32SliceVarPFlag(variable *[]float32, name string, shorthand string, value []float32, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Float32SliceVarP(variable, name, shorthand, value, usage)
	return b
}

// WithFloat64Flag defines a float64 flag with specified name, default value,
// and usage string. The return value is the address of a float64 variable that
// stores the value of the flag.
func (b *CobraCmdBuilder) WithFloat64Flag(name string, value float64, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Float64(name, value, usage)
	return b
}

// WithFloat64PFlag is like Float64, but accepts a shorthand letter that can be
// used after a single dash.
func (b *CobraCmdBuilder) WithFloat64PFlag(name string, shorthand string, value float64, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Float64P(name, shorthand, value, usage)
	return b
}

// WithFloat64VarFlag defines a float64 flag with specified name, default value,
// and usage string. The argument p points to a float64 variable in which to
// store the value of the flag.
func (b *CobraCmdBuilder) WithFloat64VarFlag(variable *float64, name string, value float64, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Float64Var(variable, name, value, usage)
	return b
}

// WithFloat64VarPFlag is like Float64Var, but accepts a shorthand letter that
// can be used after a single dash.
func (b *CobraCmdBuilder) WithFloat64VarPFlag(variable *float64, name string, shorthand string, value float64, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Float64VarP(variable, name, shorthand, value, usage)
	return b
}

// WithFloat64SliceFlag defines a []float64 flag with specified name, default
// value, and usage string. The return value is the address of a []float64
// variable that stores the value of the flag.
func (b *CobraCmdBuilder) WithFloat64SliceFlag(name string, value []float64, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Float64Slice(name, value, usage)
	return b
}

// WithFloat64SlicePFlag is like Float64Slice, but accepts a shorthand letter
// that can be used after a single dash.
func (b *CobraCmdBuilder) WithFloat64SlicePFlag(name string, shorthand string, value []float64, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Float64SliceP(name, shorthand, value, usage)
	return b
}

// WithFloat64SliceVarFlag defines a float64Slice flag with specified name,
// default value, and usage string. The argument p points to a []float64
// variable in which to store the value of the flag.
func (b *CobraCmdBuilder) WithFloat64SliceVarFlag(variable *[]float64, name string, value []float64, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Float64SliceVar(variable, name, value, usage)
	return b
}

// WithFloat64SliceVarPFlag is like Float64SliceVar, but accepts a shorthand
// letter that can be used after a single dash.
func (b *CobraCmdBuilder) WithFloat64SliceVarPFlag(variable *[]float64, name string, shorthand string, value []float64, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Float64SliceVarP(variable, name, shorthand, value, usage)
	return b
}

// WithIntFlag defines an int flag with specified name, default value, and usage
// string. The return value is the address of an int variable that stores the
// value of the flag.
func (b *CobraCmdBuilder) WithIntFlag(name string, value int, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Int(name, value, usage)
	return b
}

// WithIntPFlag is like Int, but accepts a shorthand letter that can be used
// after a single dash.
func (b *CobraCmdBuilder) WithIntPFlag(name string, shorthand string, value int, usage string) *CobraCmdBuilder {
	b.cmd.Flags().IntP(name, shorthand, value, usage)
	return b
}

// WithIntVarFlag defines an int flag with specified name, default value, and
// usage string. The argument p points to an int variable in which to store the
// value of the flag.
func (b *CobraCmdBuilder) WithIntVarFlag(variable *int, name string, value int, usage string) *CobraCmdBuilder {
	b.cmd.Flags().IntVar(variable, name, value, usage)
	return b
}

// WithIntVarPFlag is like IntVar, but accepts a shorthand letter that can be
// used after a single dash.
func (b *CobraCmdBuilder) WithIntVarPFlag(variable *int, name string, shorthand string, value int, usage string) *CobraCmdBuilder {
	b.cmd.Flags().IntVarP(variable, name, shorthand, value, usage)
	return b
}

// WithIntSliceFlag defines a []int flag with specified name, default value, and
// usage string. The return value is the address of a []int variable that stores
// the value of the flag.
func (b *CobraCmdBuilder) WithIntSliceFlag(name string, value []int, usage string) *CobraCmdBuilder {
	b.cmd.Flags().IntSlice(name, value, usage)
	return b
}

// WithIntSlicePFlag is like IntSlice, but accepts a shorthand letter that can
// be used after a single dash.
func (b *CobraCmdBuilder) WithIntSlicePFlag(name string, shorthand string, value []int, usage string) *CobraCmdBuilder {
	b.cmd.Flags().IntSliceP(name, shorthand, value, usage)
	return b
}

// WithIntSliceVarFlag defines a intSlice flag with specified name, default
// value, and usage string. The argument p points to a []int variable in which
// to store the value of the flag.
func (b *CobraCmdBuilder) WithIntSliceVarFlag(variable *[]int, name string, value []int, usage string) *CobraCmdBuilder {
	b.cmd.Flags().IntSliceVar(variable, name, value, usage)
	return b
}

// WithIntSliceVarPFlag is like IntSliceVar, but accepts a shorthand letter that
// can be used after a single dash.
func (b *CobraCmdBuilder) WithIntSliceVarPFlag(variable *[]int, name string, shorthand string, value []int, usage string) *CobraCmdBuilder {
	b.cmd.Flags().IntSliceVarP(variable, name, shorthand, value, usage)
	return b
}

// WithInt8Flag defines an int8 flag with specified name, default value, and
// usage string. The return value is the address of an int8 variable that stores
// the value of the flag.
func (b *CobraCmdBuilder) WithInt8Flag(name string, value int8, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Int8(name, value, usage)
	return b
}

// WithInt8PFlag is like Int8, but accepts a shorthand letter that can be used
// after a single dash.
func (b *CobraCmdBuilder) WithInt8PFlag(name string, shorthand string, value int8, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Int8P(name, shorthand, value, usage)
	return b
}

// WithInt8VarFlag defines an int8 flag with specified name, default value, and
// usage string. The argument p points to an int8 variable in which to store the
// value of the flag.
func (b *CobraCmdBuilder) WithInt8VarFlag(variable *int8, name string, value int8, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Int8Var(variable, name, value, usage)
	return b
}

// WithInt8VarPFlag is like Int8Var, but accepts a shorthand letter that can be
// used after a single dash.
func (b *CobraCmdBuilder) WithInt8VarPFlag(variable *int8, name string, shorthand string, value int8, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Int8VarP(variable, name, shorthand, value, usage)
	return b
}

// WithInt16Flag defines an int16 flag with specified name, default value, and
// usage string. The return value is the address of an int16 variable that
// stores the value of the flag.
func (b *CobraCmdBuilder) WithInt16Flag(name string, value int16, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Int16(name, value, usage)
	return b
}

// WithInt16PFlag is like Int16, but accepts a shorthand letter that can be used
// after a single dash.
func (b *CobraCmdBuilder) WithInt16PFlag(name string, shorthand string, value int16, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Int16P(name, shorthand, value, usage)
	return b
}

// WithInt16VarFlag defines an int16 flag with specified name, default value,
// and usage string. The argument p points to an int16 variable in which to
// store the value of the flag.
func (b *CobraCmdBuilder) WithInt16VarFlag(variable *int16, name string, value int16, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Int16Var(variable, name, value, usage)
	return b
}

// WithInt16VarPFlag is like Int16Var, but accepts a shorthand letter that can
// be used after a single dash.
func (b *CobraCmdBuilder) WithInt16VarPFlag(variable *int16, name string, shorthand string, value int16, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Int16VarP(variable, name, shorthand, value, usage)
	return b
}

// WithInt32Flag defines an int32 flag with specified name, default value, and
// usage string. The return value is the address of an int32 variable that
// stores the value of the flag.
func (b *CobraCmdBuilder) WithInt32Flag(name string, value int32, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Int32(name, value, usage)
	return b
}

// WithInt32PFlag is like Int32, but accepts a shorthand letter that can be used
// after a single dash.
func (b *CobraCmdBuilder) WithInt32PFlag(name string, shorthand string, value int32, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Int32P(name, shorthand, value, usage)
	return b
}

// WithInt32VarFlag defines an int32 flag with specified name, default value,
// and usage string. The argument p points to an int32 variable in which to
// store the value of the flag.
func (b *CobraCmdBuilder) WithInt32VarFlag(variable *int32, name string, value int32, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Int32Var(variable, name, value, usage)
	return b
}

// WithInt32VarPFlag is like Int32Var, but accepts a shorthand letter that can
// be used after a single dash.
func (b *CobraCmdBuilder) WithInt32VarPFlag(variable *int32, name string, shorthand string, value int32, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Int32VarP(variable, name, shorthand, value, usage)
	return b
}

// WithInt32SliceFlag defines a []int32 flag with specified name, default value,
// and usage string. The return value is the address of a []int32 variable that
// stores the value of the flag.
func (b *CobraCmdBuilder) WithInt32SliceFlag(name string, value []int32, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Int32Slice(name, value, usage)
	return b
}

// WithInt32SlicePFlag is like Int32Slice, but accepts a shorthand letter that
// can be used after a single dash.
func (b *CobraCmdBuilder) WithInt32SlicePFlag(name string, shorthand string, value []int32, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Int32SliceP(name, shorthand, value, usage)
	return b
}

// WithInt32SliceVarFlag defines a int32Slice flag with specified name, default
// value, and usage string. The argument p points to a []int32 variable in which
// to store the value of the flag.
func (b *CobraCmdBuilder) WithInt32SliceVarFlag(variable *[]int32, name string, value []int32, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Int32SliceVar(variable, name, value, usage)
	return b
}

// WithInt32SliceVarPFlag is like Int32SliceVar, but accepts a shorthand letter
// that can be used after a single dash.
func (b *CobraCmdBuilder) WithInt32SliceVarPFlag(variable *[]int32, name string, shorthand string, value []int32, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Int32SliceVarP(variable, name, shorthand, value, usage)
	return b
}

// WithInt64Flag defines an int64 flag with specified name, default value, and
// usage string. The return value is the address of an int64 variable that
// stores the value of the flag.
func (b *CobraCmdBuilder) WithInt64Flag(name string, value int64, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Int64(name, value, usage)
	return b
}

// WithInt64PFlag is like Int64, but accepts a shorthand letter that can be used
// after a single dash.
func (b *CobraCmdBuilder) WithInt64PFlag(name string, shorthand string, value int64, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Int64P(name, shorthand, value, usage)
	return b
}

// WithInt64VarFlag defines an int64 flag with specified name, default value,
// and usage string. The argument p points to an int64 variable in which to
// store the value of the flag.
func (b *CobraCmdBuilder) WithInt64VarFlag(variable *int64, name string, value int64, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Int64Var(variable, name, value, usage)
	return b
}

// WithInt64VarPFlag is like Int64Var, but accepts a shorthand letter that can
// be used after a single dash.
func (b *CobraCmdBuilder) WithInt64VarPFlag(variable *int64, name string, shorthand string, value int64, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Int64VarP(variable, name, shorthand, value, usage)
	return b
}

// WithInt64SliceFlag defines a []int64 flag with specified name, default value,
// and usage string. The return value is the address of a []int64 variable that
// stores the value of the flag.
func (b *CobraCmdBuilder) WithInt64SliceFlag(name string, value []int64, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Int64Slice(name, value, usage)
	return b
}

// WithInt64SlicePFlag is like Int64Slice, but accepts a shorthand letter that
// can be used after a single dash.
func (b *CobraCmdBuilder) WithInt64SlicePFlag(name string, shorthand string, value []int64, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Int64SliceP(name, shorthand, value, usage)
	return b
}

// WithInt64SliceVarFlag defines a int64Slice flag with specified name, default
// value, and usage string. The argument p points to a []int64 variable in which
// to store the value of the flag.
func (b *CobraCmdBuilder) WithInt64SliceVarFlag(variable *[]int64, name string, value []int64, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Int64SliceVar(variable, name, value, usage)
	return b
}

// WithInt64SliceVarPFlag is like Int64SliceVar, but accepts a shorthand letter
// that can be used after a single dash.
func (b *CobraCmdBuilder) WithInt64SliceVarPFlag(variable *[]int64, name string, shorthand string, value []int64, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Int64SliceVarP(variable, name, shorthand, value, usage)
	return b
}

// WithUintFlag defines a uint flag with specified name, default value, and
// usage string. The return value is the address of a uint variable that stores
// the value of the flag.
func (b *CobraCmdBuilder) WithUintFlag(name string, value uint, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Uint(name, value, usage)
	return b
}

// WithIntPFlag is like Uint, but accepts a shorthand letter that can be used
// after a single dash.
func (b *CobraCmdBuilder) WithUintPFlag(name string, shorthand string, value uint, usage string) *CobraCmdBuilder {
	b.cmd.Flags().UintP(name, shorthand, value, usage)
	return b
}

// WithUintVarFlag defines a uint flag with specified name, default value, and
// usage string. The argument p points to a uint variable in which to store the
// value of the flag.
func (b *CobraCmdBuilder) WithUintVarFlag(variable *uint, name string, value uint, usage string) *CobraCmdBuilder {
	b.cmd.Flags().UintVar(variable, name, value, usage)
	return b
}

// WithUintVarPFlag is like UintVar, but accepts a shorthand letter that can be
// used after a single dash.
func (b *CobraCmdBuilder) WithUintVarPFlag(variable *uint, name string, shorthand string, value uint, usage string) *CobraCmdBuilder {
	b.cmd.Flags().UintVarP(variable, name, shorthand, value, usage)
	return b
}

// WithUintSliceFlag defines a []uint flag with specified name, default value,
// and usage string. The return value is the address of a []uint variable that
// stores the value of the flag.
func (b *CobraCmdBuilder) WithUintSliceFlag(name string, value []uint, usage string) *CobraCmdBuilder {
	b.cmd.Flags().UintSlice(name, value, usage)
	return b
}

// WithUintSlicePFlag is like UintSlice, but accepts a shorthand letter that can
// be used after a single dash.
func (b *CobraCmdBuilder) WithUintSlicePFlag(name string, shorthand string, value []uint, usage string) *CobraCmdBuilder {
	b.cmd.Flags().UintSliceP(name, shorthand, value, usage)
	return b
}

// WithUintSliceVarFlag defines a uintSlice flag with specified name, default
// value, and usage string. The argument p points to a []uint variable in which
// to store the value of the flag.
func (b *CobraCmdBuilder) WithUintSliceVarFlag(variable *[]uint, name string, value []uint, usage string) *CobraCmdBuilder {
	b.cmd.Flags().UintSliceVar(variable, name, value, usage)
	return b
}

// WithUintSliceVarPFlag is like UintSliceVar, but accepts a shorthand letter
// that can be used after a single dash.
func (b *CobraCmdBuilder) WithUintSliceVarPFlag(variable *[]uint, name string, shorthand string, value []uint, usage string) *CobraCmdBuilder {
	b.cmd.Flags().UintSliceVarP(variable, name, shorthand, value, usage)
	return b
}

// WithUint8Flag defines a uint8 flag with specified name, default value, and
// usage string. The return value is the address of a uint8 variable that stores
// the value of the flag.
func (b *CobraCmdBuilder) WithUint8Flag(name string, value uint8, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Uint8(name, value, usage)
	return b
}

// WithUint8PFlag is like Uint8, but accepts a shorthand letter that can be used
// after a single dash.
func (b *CobraCmdBuilder) WithUint8PFlag(name string, shorthand string, value uint8, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Uint8P(name, shorthand, value, usage)
	return b
}

// WithUint8VarFlag defines a uint8 flag with specified name, default value, and
// usage string. The argument p points to a uint8 variable in which to store the
// value of the flag.
func (b *CobraCmdBuilder) WithUint8VarFlag(variable *uint8, name string, value uint8, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Uint8Var(variable, name, value, usage)
	return b
}

// WithUint8VarPFlag is like Uint8Var, but accepts a shorthand letter that can
// be used after a single dash.
func (b *CobraCmdBuilder) WithUint8VarPFlag(variable *uint8, name string, shorthand string, value uint8, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Uint8VarP(variable, name, shorthand, value, usage)
	return b
}

// WithUint16Flag defines a uint flag with specified name, default value, and
// usage string. The return value is the address of a uint variable that stores
// the value of the flag.
func (b *CobraCmdBuilder) WithUint16Flag(name string, value uint16, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Uint16(name, value, usage)
	return b
}

// WithUint16PFlag is like Uint16, but accepts a shorthand letter that can be
// used after a single dash.
func (b *CobraCmdBuilder) WithUint16PFlag(name string, shorthand string, value uint16, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Uint16P(name, shorthand, value, usage)
	return b
}

// WithUint16VarFlag defines a uint flag with specified name, default value, and
// usage string. The argument p points to a uint variable in which to store the
// value of the flag.
func (b *CobraCmdBuilder) WithUint16VarFlag(variable *uint16, name string, value uint16, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Uint16Var(variable, name, value, usage)
	return b
}

// WithUint16VarPFlag is like Uint16Var, but accepts a shorthand letter that can
// be used after a single dash.
func (b *CobraCmdBuilder) WithUint16VarPFlag(variable *uint16, name string, shorthand string, value uint16, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Uint16VarP(variable, name, shorthand, value, usage)
	return b
}

// WithUint32Flag defines a uint32 flag with specified name, default value, and
// usage string. The return value is the address of a uint32 variable that
// stores the value of the flag.
func (b *CobraCmdBuilder) WithUint32Flag(name string, value uint32, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Uint32(name, value, usage)
	return b
}

// WithUint32PFlag is like Uint32, but accepts a shorthand letter that can be
// used after a single dash.
func (b *CobraCmdBuilder) WithUint32PFlag(name string, shorthand string, value uint32, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Uint32P(name, shorthand, value, usage)
	return b
}

// WithUint32VarFlag defines a uint32 flag with specified name, default value,
// and usage string. The argument p points to a uint32 variable in which to
// store the value of the flag.
func (b *CobraCmdBuilder) WithUint32VarFlag(variable *uint32, name string, value uint32, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Uint32Var(variable, name, value, usage)
	return b
}

// WithUint32VarPFlag is like Uint32Var, but accepts a shorthand letter that can
// be used after a single dash.
func (b *CobraCmdBuilder) WithUint32VarPFlag(variable *uint32, name string, shorthand string, value uint32, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Uint32VarP(variable, name, shorthand, value, usage)
	return b
}

// WithUint64Flag defines a uint64 flag with specified name, default value, and
// usage string. The return value is the address of a uint64 variable that
// stores the value of the flag.
func (b *CobraCmdBuilder) WithUint64Flag(name string, value uint64, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Uint64(name, value, usage)
	return b
}

// WithUint64PFlag is like Uint64, but accepts a shorthand letter that can be
// used after a single dash.
func (b *CobraCmdBuilder) WithUint64PFlag(name string, shorthand string, value uint64, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Uint64P(name, shorthand, value, usage)
	return b
}

// WithUint64VarFlag defines a uint64 flag with specified name, default value,
// and usage string. The argument p points to a uint64 variable in which to
// store the value of the flag.
func (b *CobraCmdBuilder) WithUint64VarFlag(variable *uint64, name string, value uint64, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Uint64Var(variable, name, value, usage)
	return b
}

// WithUint64VarPFlag is like Uint64Var, but accepts a shorthand letter that can
// be used after a single dash.
func (b *CobraCmdBuilder) WithUint64VarPFlag(variable *uint64, name string, shorthand string, value uint64, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Uint64VarP(variable, name, shorthand, value, usage)
	return b
}

// WithStringFlag defines a string flag with specified name, default value, and
// usage string. The return value is the address of a string variable that
// stores the value of the flag.
func (b *CobraCmdBuilder) WithStringFlag(name string, value string, usage string) *CobraCmdBuilder {
	b.cmd.Flags().String(name, value, usage)
	return b
}

// WithStringPFlag is like String, but accepts a shorthand letter that can be
// used after a single dash.
func (b *CobraCmdBuilder) WithStringPFlag(name string, shorthand string, value string, usage string) *CobraCmdBuilder {
	b.cmd.Flags().StringP(name, shorthand, value, usage)
	return b
}

// WithStringVarFlag defines a string flag with specified name, default value,
// and usage string. The argument p points to a string variable in which to
// store the value of the flag.
func (b *CobraCmdBuilder) WithStringVarFlag(variable *string, name string, value string, usage string) *CobraCmdBuilder {
	b.cmd.Flags().StringVar(variable, name, value, usage)
	return b
}

// WithStringVarPFlag is like StringVar, but accepts a shorthand letter that can
// be used after a single dash.
func (b *CobraCmdBuilder) WithStringVarPFlag(variable *string, name string, shorthand string, value string, usage string) *CobraCmdBuilder {
	b.cmd.Flags().StringVarP(variable, name, shorthand, value, usage)
	return b
}

// WithStringSliceFlag defines a string flag with specified name, default value,
// and usage string. The return value is the address of a []string variable that
// stores the value of the flag. Compared to StringArray flags, StringSlice
// flags take comma-separated value as arguments and split them accordingly.
// For example:
//
//	--ss="v1,v2" --ss="v3"
//
// will result in
//
//	[]string{"v1", "v2", "v3"}
func (b *CobraCmdBuilder) WithStringSliceFlag(name string, value []string, usage string) *CobraCmdBuilder {
	b.cmd.Flags().StringSlice(name, value, usage)
	return b
}

// WithStringSlicePFlag is like StringSlice, but accepts a shorthand letter that
// can be used after a single dash.
func (b *CobraCmdBuilder) WithStringSlicePFlag(name string, shorthand string, value []string, usage string) *CobraCmdBuilder {
	b.cmd.Flags().StringSliceP(name, shorthand, value, usage)
	return b
}

// WithStringSliceVarFlag defines a string flag with specified name, default
// value, and usage string. The argument p points to a []string variable in
// which to store the value of the flag. Compared to StringArray flags,
// StringSlice flags take comma-separated value as arguments and split them
// accordingly.
// For example:
//
//	--ss="v1,v2" --ss="v3"
//
// will result in
//
//	[]string{"v1", "v2", "v3"}
func (b *CobraCmdBuilder) WithStringSliceVarFlag(variable *[]string, name string, value []string, usage string) *CobraCmdBuilder {
	b.cmd.Flags().StringSliceVar(variable, name, value, usage)
	return b
}

// WithStringSliceVarPFlag is like StringSliceVar, but accepts a shorthand
// letter that can be used after a single dash.
func (b *CobraCmdBuilder) WithStringSliceVarPFlag(variable *[]string, name string, shorthand string, value []string, usage string) *CobraCmdBuilder {
	b.cmd.Flags().StringSliceVarP(variable, name, shorthand, value, usage)
	return b
}

// WithStringArrayFlag defines a string flag with specified name, default value,
// and usage string. The return value is the address of a []string variable that
// stores the value of the flag. The value of each argument will not try to be
// separated by comma. Use a StringSlice for that.
func (b *CobraCmdBuilder) WithStringArrayFlag(name string, value []string, usage string) *CobraCmdBuilder {
	b.cmd.Flags().StringArray(name, value, usage)
	return b
}

// WithStringArrayPFlag is like StringArray, but accepts a shorthand letter that
// can be used after a single dash.
func (b *CobraCmdBuilder) WithStringArrayPFlag(name string, shorthand string, value []string, usage string) *CobraCmdBuilder {
	b.cmd.Flags().StringArrayP(name, shorthand, value, usage)
	return b
}

// WithStringArrayVarFlag defines a string flag with specified name, default
// value, and usage string. The argument p points to a []string variable in
// which to store the values of the multiple flags. The value of each argument
// will not try to be separated by comma. Use a StringSlice for that.
func (b *CobraCmdBuilder) WithStringArrayVarFlag(variable *[]string, name string, value []string, usage string) *CobraCmdBuilder {
	b.cmd.Flags().StringArrayVar(variable, name, value, usage)
	return b
}

// WithStringArrayVarPFlag is like StringArrayVar, but accepts a shorthand
// letter that can be used after a single dash.
func (b *CobraCmdBuilder) WithStringArrayVarPFlag(variable *[]string, name string, shorthand string, value []string, usage string) *CobraCmdBuilder {
	b.cmd.Flags().StringArrayVarP(variable, name, shorthand, value, usage)
	return b
}

// WithStringToIntFlag defines a string flag with specified name, default value,
// and usage string. The return value is the address of a map[string]int
// variable that stores the value of the flag. The value of each argument will
// not try to be separated by comma
func (b *CobraCmdBuilder) WithStringToIntFlag(name string, value map[string]int, usage string) *CobraCmdBuilder {
	b.cmd.Flags().StringToInt(name, value, usage)
	return b
}

// WithStringToIntPFlag is like StringToInt, but accepts a shorthand letter that
// can be used after a single dash.
func (b *CobraCmdBuilder) WithStringToIntPFlag(name string, shorthand string, value map[string]int, usage string) *CobraCmdBuilder {
	b.cmd.Flags().StringToIntP(name, shorthand, value, usage)
	return b
}

// WithStringToIntVarFlag defines a string flag with specified name, default
// value, and usage string. The argument p points to a map[string]int variable
// in which to store the values of the multiple flags. The value of each
// argument will not try to be separated by comma
func (b *CobraCmdBuilder) WithStringToIntVarFlag(variable *map[string]int, name string, value map[string]int, usage string) *CobraCmdBuilder {
	b.cmd.Flags().StringToIntVar(variable, name, value, usage)
	return b
}

// WithStringToIntVarPFlag is like StringToIntVar, but accepts a shorthand
// letter that can be used after a single dash.
func (b *CobraCmdBuilder) WithStringToIntVarPFlag(variable *map[string]int, name string, shorthand string, value map[string]int, usage string) *CobraCmdBuilder {
	b.cmd.Flags().StringToIntVarP(variable, name, shorthand, value, usage)
	return b
}

// WithStringToInt64Flag defines a string flag with specified name, default
// value, and usage string. The return value is the address of a map[string
// ]int64 variable that stores the value of the flag. The value of each argument
// will not try to be separated by comma
func (b *CobraCmdBuilder) WithStringToInt64Flag(name string, value map[string]int64, usage string) *CobraCmdBuilder {
	b.cmd.Flags().StringToInt64(name, value, usage)
	return b
}

// WithStringToInt64PFlag is like StringToInt64, but accepts a shorthand letter
// that can be used after a single dash.
func (b *CobraCmdBuilder) WithStringToInt64PFlag(name string, shorthand string, value map[string]int64, usage string) *CobraCmdBuilder {
	b.cmd.Flags().StringToInt64P(name, shorthand, value, usage)
	return b
}

// WithStringToInt64VarFlag defines a string flag with specified name, default
// value, and usage string. The argument p point64s to a map[string]int64
// variable in which to store the values of the multiple flags. The value of
// each argument will not try to be separated by comma
func (b *CobraCmdBuilder) WithStringToInt64VarFlag(variable *map[string]int64, name string, value map[string]int64, usage string) *CobraCmdBuilder {
	b.cmd.Flags().StringToInt64Var(variable, name, value, usage)
	return b
}

// WithStringToInt64VarPFlag is like StringToInt64Var, but accepts a shorthand
// letter that can be used after a single dash.
func (b *CobraCmdBuilder) WithStringToInt64VarPFlag(variable *map[string]int64, name string, shorthand string, value map[string]int64, usage string) *CobraCmdBuilder {
	b.cmd.Flags().StringToInt64VarP(variable, name, shorthand, value, usage)
	return b
}

// WithIPFlag defines an net.IP flag with specified name, default value, and
// usage string. The return value is the address of an net.IP variable that
// stores the value of the flag.
func (b *CobraCmdBuilder) WithIPFlag(name string, value net.IP, usage string) *CobraCmdBuilder {
	b.cmd.Flags().IP(name, value, usage)
	return b
}

// WithIPPFlag is like IP, but accepts a shorthand letter that can be used after
// a single dash.
func (b *CobraCmdBuilder) WithIPPFlag(name string, shorthand string, value net.IP, usage string) *CobraCmdBuilder {
	b.cmd.Flags().IPP(name, shorthand, value, usage)
	return b
}

// WithIPVarFlag defines an net.IP flag with specified name, default value, and
// usage string. The argument p points to an net.IP variable in which to store
// the value of the flag.
func (b *CobraCmdBuilder) WithIPVarFlag(variable *net.IP, name string, value net.IP, usage string) *CobraCmdBuilder {
	b.cmd.Flags().IPVar(variable, name, value, usage)
	return b
}

// WithIPVarPFlag is like IPVar, but accepts a shorthand letter that can be used
// after a single dash.
func (b *CobraCmdBuilder) WithIPVarPFlag(variable *net.IP, name string, shorthand string, value net.IP, usage string) *CobraCmdBuilder {
	b.cmd.Flags().IPVarP(variable, name, shorthand, value, usage)
	return b
}

// WithIPSliceFlag defines a []net.IP flag with specified name, default value,
// and usage string. The return value is the address of a []net.IP variable that
// stores the value of that flag.
func (b *CobraCmdBuilder) WithIPSliceFlag(name string, value []net.IP, usage string) *CobraCmdBuilder {
	b.cmd.Flags().IPSlice(name, value, usage)
	return b
}

// WithIPSlicePFlag is like IPSlice, but accepts a shorthand letter that can be
// used after a single dash.
func (b *CobraCmdBuilder) WithIPSlicePFlag(name string, shorthand string, value []net.IP, usage string) *CobraCmdBuilder {
	b.cmd.Flags().IPSliceP(name, shorthand, value, usage)
	return b
}

// WithIPSliceVarFlag defines a ipSlice flag with specified name, default value,
// and usage string. The argument p points to a []net.IP variable in which to
// store the value of the flag.
func (b *CobraCmdBuilder) WithIPSliceVarFlag(variable *[]net.IP, name string, value []net.IP, usage string) *CobraCmdBuilder {
	b.cmd.Flags().IPSliceVar(variable, name, value, usage)
	return b
}

// WithIPSliceVarPFlag is like IPSliceVar, but accepts a shorthand letter that
// can be used after a single dash.
func (b *CobraCmdBuilder) WithIPSliceVarPFlag(variable *[]net.IP, name string, shorthand string, value []net.IP, usage string) *CobraCmdBuilder {
	b.cmd.Flags().IPSliceVarP(variable, name, shorthand, value, usage)
	return b
}

// WithIPMaskFlag defines an net.IPMask flag with specified name, default value,
// and usage string. The return value is the address of an net.IPMask variable
// that stores the value of the flag.
func (b *CobraCmdBuilder) WithIPMaskFlag(name string, value net.IPMask, usage string) *CobraCmdBuilder {
	b.cmd.Flags().IPMask(name, value, usage)
	return b
}

// WithIPMaskPFlag is like IPMask, but accepts a shorthand letter that can be
// used after a single dash.
func (b *CobraCmdBuilder) WithIPMaskPFlag(name string, shorthand string, value net.IPMask, usage string) *CobraCmdBuilder {
	b.cmd.Flags().IPMaskP(name, shorthand, value, usage)
	return b
}

// WithIPMaskVarFlag defines an net.IPMask flag with specified name, default
// value, and usage string. The argument p points to an net.IPMask variable in
// which to store the value of the flag.
func (b *CobraCmdBuilder) WithIPMaskVarFlag(variable *net.IPMask, name string, value net.IPMask, usage string) *CobraCmdBuilder {
	b.cmd.Flags().IPMaskVar(variable, name, value, usage)
	return b
}

// WithIPMaskVarPFlag is like IPMaskVar, but accepts a shorthand letter that can
// be used after a single dash.
func (b *CobraCmdBuilder) WithIPMaskVarPFlag(variable *net.IPMask, name string, shorthand string, value net.IPMask, usage string) *CobraCmdBuilder {
	b.cmd.Flags().IPMaskVarP(variable, name, shorthand, value, usage)
	return b
}

// WithIPNetFlag defines an net.IPNet flag with specified name, default value,
// and usage string. The return value is the address of an net.IPNet variable
// that stores the value of the flag.
func (b *CobraCmdBuilder) WithIPNetFlag(name string, value net.IPNet, usage string) *CobraCmdBuilder {
	b.cmd.Flags().IPNet(name, value, usage)
	return b
}

// WithIPNetPFlag is like IPNet, but accepts a shorthand letter that can be used
// after a single dash.
func (b *CobraCmdBuilder) WithIPNetPFlag(name string, shorthand string, value net.IPNet, usage string) *CobraCmdBuilder {
	b.cmd.Flags().IPNetP(name, shorthand, value, usage)
	return b
}

// WithIPNetVarFlag defines an net.IPNet flag with specified name, default
// value, and usage string. The argument p points to an net.IPNet variable in
// which to store the value of the flag.
func (b *CobraCmdBuilder) WithIPNetVarFlag(variable *net.IPNet, name string, value net.IPNet, usage string) *CobraCmdBuilder {
	b.cmd.Flags().IPNetVar(variable, name, value, usage)
	return b
}

// WithIPNetVarPFlag is like IPNetVar, but accepts a shorthand letter that can
// be used after a single dash.
func (b *CobraCmdBuilder) WithIPNetVarPFlag(variable *net.IPNet, name string, shorthand string, value net.IPNet, usage string) *CobraCmdBuilder {
	b.cmd.Flags().IPNetVarP(variable, name, shorthand, value, usage)
	return b
}

// WithVarFlag defines a flag with the specified name and usage string. The type
// and value of the flag are represented by the first argument, of type Value,
// which typically holds a user-defined implementation of Value. For instance,
// the caller could create a flag that turns a comma-separated string into a
// slice of strings by giving the slice the methods of Value; in particular, Set
// would decompose the comma-separated string into the slice.
func (b *CobraCmdBuilder) WithVarFlag(value pflag.Value, name string, usage string) *CobraCmdBuilder {
	b.cmd.Flags().Var(value, name, usage)
	return b
}

// WithVarPFlag is like Var, but accepts a shorthand letter that can be used
// after a single dash.
func (b *CobraCmdBuilder) WithVarPFlag(value pflag.Value, name string, shorthand string, usage string) *CobraCmdBuilder {
	b.cmd.Flags().VarP(value, name, shorthand, usage)
	return b
}

// MarkFlagHidden sets a flag to 'hidden' in your program. It will continue to
// function but will not show up in help or usage messages.
func (b *CobraCmdBuilder) MarkFlagHidden(name string) *CobraCmdBuilder {
	err := b.cmd.Flags().MarkHidden(name)
	if err != nil {
		panic(err)
	}
	return b
}

// MarkFlagDeprecated indicated that a flag is deprecated in your program. It
// will continue to function but will not show up in help or usage messages.
// Using this flag will also print the given usageMessage.
func (b *CobraCmdBuilder) MarkFlagDeprecated(name string, usage string) *CobraCmdBuilder {
	err := b.cmd.Flags().MarkDeprecated(name, usage)
	if err != nil {
		panic(err)
	}
	return b
}

// MarkFlagShorthandDeprecated will mark the shorthand of a flag deprecated in
// your program. It will continue to function but will not show up in help or
// usage messages. Using this flag will also print the given usageMessage.
func (b *CobraCmdBuilder) MarkFlagShorthandDeprecated(name string, usage string) *CobraCmdBuilder {
	err := b.cmd.Flags().MarkShorthandDeprecated(name, usage)
	if err != nil {
		panic(err)
	}
	return b
}

// WithBoolPersistentFlag defines a bool flag with specified name, default
// value, and usage string. The return value is the address of a bool variable
// that stores the value of the flag.
func (b *CobraCmdBuilder) WithBoolPersistentFlag(name string, value bool, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Bool(name, value, usage)
	return b
}

// WithBoolPPersistentFlag BoolP is like Bool, but accepts a shorthand letter
// that can be used after a single dash.
func (b *CobraCmdBuilder) WithBoolPPersistentFlag(name string, shorthand string, value bool, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().BoolP(name, shorthand, value, usage)
	return b
}

// WithBoolVarPersistentFlag defines a bool flag with specified name, default
// value, and usage string. The argument p points to a bool variable in which to
// store the value of the flag.
func (b *CobraCmdBuilder) WithBoolVarPersistentFlag(variable *bool, name string, value bool, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().BoolVar(variable, name, value, usage)
	return b
}

// WithBoolVarPPersistentFlag is like BoolVar, but accepts a shorthand letter
// that can be used after a single dash.
func (b *CobraCmdBuilder) WithBoolVarPPersistentFlag(variable *bool, name string, shorthand string, value bool, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().BoolVarP(variable, name, shorthand, value, usage)
	return b
}

// WithBoolSlicePersistentFlag defines a []bool flag with specified name,
// default value, and usage string. The return value is the address of a []bool
// variable that stores the value of the flag.
func (b *CobraCmdBuilder) WithBoolSlicePersistentFlag(name string, value []bool, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().BoolSlice(name, value, usage)
	return b
}

// WithBoolSlicePPersistentFlag is like BoolSlice, but accepts a shorthand
// letter that can be used after a single dash.
func (b *CobraCmdBuilder) WithBoolSlicePPersistentFlag(name string, shorthand string, value []bool, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().BoolSliceP(name, shorthand, value, usage)
	return b
}

// WithBoolSliceVarPersistentFlag defines a boolSlice flag with specified name,
// default value, and usage string. The argument p points to a []bool variable
// in which to store the value of the flag.
func (b *CobraCmdBuilder) WithBoolSliceVarPersistentFlag(variable *[]bool, name string, value []bool, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().BoolSliceVar(variable, name, value, usage)
	return b
}

// WithBoolSliceVarPPersistentFlag is like BoolSliceVar, but accepts a shorthand
// letter that can be used after a single dash.
func (b *CobraCmdBuilder) WithBoolSliceVarPPersistentFlag(variable *[]bool, name string, shorthand string, value []bool, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().BoolSliceVarP(variable, name, shorthand, value, usage)
	return b
}

// WithBytesBase64PersistentFlag defines an []byte flag with specified name,
// default value, and usage string. The return value is the address of an []byte
// variable that stores the value of the flag.
func (b *CobraCmdBuilder) WithBytesBase64PersistentFlag(name string, value []byte, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().BytesBase64(name, value, usage)
	return b
}

// WithBytesBase64PPersistentFlag is like BytesBase64, but accepts a shorthand
// letter that can be used after a single dash.
func (b *CobraCmdBuilder) WithBytesBase64PPersistentFlag(name string, shorthand string, value []byte, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().BytesBase64P(name, shorthand, value, usage)
	return b
}

// WithBytesBase64VarPersistentFlag defines an []byte flag with specified name,
// default value, and usage string. The argument p points to an []byte variable
// in which to store the value of the flag.
func (b *CobraCmdBuilder) WithBytesBase64VarPersistentFlag(variable *[]byte, name string, value []byte, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().BytesBase64Var(variable, name, value, usage)
	return b
}

// WithBytesBase64VarPPersistentFlag BytesBase64VarP is like BytesBase64Var, but
// accepts a shorthand letter that can be used after a single dash.
func (b *CobraCmdBuilder) WithBytesBase64VarPPersistentFlag(variable *[]byte, name string, shorthand string, value []byte, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().BytesBase64VarP(variable, name, shorthand, value, usage)
	return b
}

// WithBytesHexPersistentFlag defines an []byte flag with specified name,
// default value, and usage string. The return value is the address of an []byte
// variable that stores the value of the flag.
func (b *CobraCmdBuilder) WithBytesHexPersistentFlag(name string, value []byte, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().BytesHex(name, value, usage)
	return b
}

// WithBytesHexPPersistentFlag is like BytesHex, but accepts a shorthand letter
// that can be used after a single dash.
func (b *CobraCmdBuilder) WithBytesHexPPersistentFlag(name string, shorthand string, value []byte, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().BytesHexP(name, shorthand, value, usage)
	return b
}

// WithBytesHexVarPersistentFlag BytesHexVar defines an []byte flag with
// specified name, default value, and usage string. The argument p points to an
// []byte variable in which to store the value of the flag.
func (b *CobraCmdBuilder) WithBytesHexVarPersistentFlag(variable *[]byte, name string, value []byte, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().BytesHexVar(variable, name, value, usage)
	return b
}

// WithBytesHexVarPPersistentFlag is like BytesHexVar, but accepts a shorthand
// letter that can be used after a single dash.
func (b *CobraCmdBuilder) WithBytesHexVarPPersistentFlag(variable *[]byte, name string, shorthand string, value []byte, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().BytesHexVarP(variable, name, shorthand, value, usage)
	return b
}

// WithCountPersistentFlag defines a count flag with specified name, default
// value, and usage string. The return value is the address of an int variable
// that stores the value of the flag. A count flag will add 1 to its value every
// time it is found on the command line.
func (b *CobraCmdBuilder) WithCountPersistentFlag(name string, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Count(name, usage)
	return b
}

// WithCountPPersistentFlag is like Count only takes a shorthand for the flag
// name.
func (b *CobraCmdBuilder) WithCountPPersistentFlag(name string, shorthand string, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().CountP(name, shorthand, usage)
	return b
}

// WithCountVarPersistentFlag defines a count flag with specified name, default
// value, and usage string. The argument p points to an int variable in which to
// store the value of the flag. A count flag will add 1 to its value every time
// it is found on the command line
func (b *CobraCmdBuilder) WithCountVarPersistentFlag(variable *int, name string, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().CountVar(variable, name, usage)
	return b
}

// WithCountVarPPersistentFlag is like CountVar only take a shorthand for the
// flag name.
func (b *CobraCmdBuilder) WithCountVarPPersistentFlag(variable *int, name string, shorthand string, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().CountVarP(variable, name, shorthand, usage)
	return b
}

// WithDurationPersistentFlag defines a time.Duration flag with specified name,
// default value, and usage string. The return value is the address of a time
// .Duration variable that stores the value of the flag.
func (b *CobraCmdBuilder) WithDurationPersistentFlag(name string, value time.Duration, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Duration(name, value, usage)
	return b
}

// WithDurationPPersistentFlag is like Duration, but accepts a shorthand letter
// that can be used after a single dash.
func (b *CobraCmdBuilder) WithDurationPPersistentFlag(name string, shorthand string, value time.Duration, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().DurationP(name, shorthand, value, usage)
	return b
}

// WithDurationVarPersistentFlag defines a time.Duration flag with specified
// name, default value, and usage string. The argument p points to a time
// .Duration variable in which to store the value of the flag.
func (b *CobraCmdBuilder) WithDurationVarPersistentFlag(variable *time.Duration, name string, value time.Duration, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().DurationVar(variable, name, value, usage)
	return b
}

// WithDurationVarPPersistentFlag is like DurationVar, but accepts a shorthand
// letter that can be used after a single dash.
func (b *CobraCmdBuilder) WithDurationVarPPersistentFlag(variable *time.Duration, name string, shorthand string, value time.Duration, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().DurationVarP(variable, name, shorthand, value, usage)
	return b
}

// WithDurationSlicePersistentFlag defines a []time.Duration flag with specified
// name, default value, and usage string. The return value is the address of a
// []time.Duration variable that stores the value of the flag.
func (b *CobraCmdBuilder) WithDurationSlicePersistentFlag(name string, value []time.Duration, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().DurationSlice(name, value, usage)
	return b
}

// WithDurationSlicePPersistentFlag is like DurationSlice, but accepts a
// shorthand letter that can be used after a single dash.
func (b *CobraCmdBuilder) WithDurationSlicePPersistentFlag(name string, shorthand string, value []time.Duration, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().DurationSliceP(name, shorthand, value, usage)
	return b
}

// WithDurationSliceVarPersistentFlag defines a durationSlice flag with
// specified name, default value, and usage string. The argument p points to a
// []time.Duration variable in which to store the value of the flag.
func (b *CobraCmdBuilder) WithDurationSliceVarPersistentFlag(variable *[]time.Duration, name string, value []time.Duration, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().DurationSliceVar(variable, name, value, usage)
	return b
}

// WithDurationSliceVarPPersistentFlag is like DurationSliceVar, but accepts a
// shorthand letter that can be used after a single dash.
func (b *CobraCmdBuilder) WithDurationSliceVarPPersistentFlag(variable *[]time.Duration, name string, shorthand string, value []time.Duration, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().DurationSliceVarP(variable, name, shorthand, value, usage)
	return b
}

// WithFloat32PersistentFlag defines a float32 flag with specified name, default
// value, and usage string. The return value is the address of a float32
// variable that stores the value of the flag.
func (b *CobraCmdBuilder) WithFloat32PersistentFlag(name string, value float32, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Float32(name, value, usage)
	return b
}

// WithFloat32PPersistentFlag is like Float32, but accepts a shorthand letter
// that can be used after a single dash.
func (b *CobraCmdBuilder) WithFloat32PPersistentFlag(name string, shorthand string, value float32, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Float32P(name, shorthand, value, usage)
	return b
}

// WithFloat32VarPersistentFlag defines a float32 flag with specified name,
// default value, and usage string. The argument p points to a float32 variable
// in which to store the value of the flag.
func (b *CobraCmdBuilder) WithFloat32VarPersistentFlag(variable *float32, name string, value float32, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Float32Var(variable, name, value, usage)
	return b
}

// WithFloat32VarPPersistentFlag is like Float32Var, but accepts a shorthand
// letter that can be used after a single dash.
func (b *CobraCmdBuilder) WithFloat32VarPPersistentFlag(variable *float32, name string, shorthand string, value float32, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Float32VarP(variable, name, shorthand, value, usage)
	return b
}

// WithFloat32SlicePersistentFlag defines a []float32 flag with specified name,
// default value, and usage string. The return value is the address of a [
// ]float32 variable that stores the value of the flag.
func (b *CobraCmdBuilder) WithFloat32SlicePersistentFlag(name string, value []float32, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Float32Slice(name, value, usage)
	return b
}

// WithFloat32SlicePPersistentFlag is like Float32Slice, but accepts a shorthand
// letter that can be used after a single dash.
func (b *CobraCmdBuilder) WithFloat32SlicePPersistentFlag(name string, shorthand string, value []float32, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Float32SliceP(name, shorthand, value, usage)
	return b
}

// WithFloat32SliceVarPersistentFlag defines a float32Slice flag with specified
// name, default value, and usage string. The argument p points to a []float32
// variable in which to store the value of the flag.
func (b *CobraCmdBuilder) WithFloat32SliceVarPersistentFlag(variable *[]float32, name string, value []float32, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Float32SliceVar(variable, name, value, usage)
	return b
}

// WithFloat32SliceVarPPersistentFlag is like Float32SliceVar, but accepts a
// shorthand letter that can be used after a single dash.
func (b *CobraCmdBuilder) WithFloat32SliceVarPPersistentFlag(variable *[]float32, name string, shorthand string, value []float32, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Float32SliceVarP(variable, name, shorthand, value, usage)
	return b
}

// WithFloat64PersistentFlag defines a float64 flag with specified name, default
// value, and usage string. The return value is the address of a float64
// variable that stores the value of the flag.
func (b *CobraCmdBuilder) WithFloat64PersistentFlag(name string, value float64, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Float64(name, value, usage)
	return b
}

// WithFloat64PPersistentFlag is like Float64, but accepts a shorthand letter
// that can be used after a single dash.
func (b *CobraCmdBuilder) WithFloat64PPersistentFlag(name string, shorthand string, value float64, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Float64P(name, shorthand, value, usage)
	return b
}

// WithFloat64VarPersistentFlag defines a float64 flag with specified name,
// default value, and usage string. The argument p points to a float64 variable
// in which to store the value of the flag.
func (b *CobraCmdBuilder) WithFloat64VarPersistentFlag(variable *float64, name string, value float64, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Float64Var(variable, name, value, usage)
	return b
}

// WithFloat64VarPPersistentFlag is like Float64Var, but accepts a shorthand
// letter that can be used after a single dash.
func (b *CobraCmdBuilder) WithFloat64VarPPersistentFlag(variable *float64, name string, shorthand string, value float64, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Float64VarP(variable, name, shorthand, value, usage)
	return b
}

// WithFloat64SlicePersistentFlag defines a []float64 flag with specified name,
// default value, and usage string. The return value is the address of a
// []float64 variable that stores the value of the flag.
func (b *CobraCmdBuilder) WithFloat64SlicePersistentFlag(name string, value []float64, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Float64Slice(name, value, usage)
	return b
}

// WithFloat64SlicePPersistentFlag is like Float64Slice, but accepts a shorthand
// letter that can be used after a single dash.
func (b *CobraCmdBuilder) WithFloat64SlicePPersistentFlag(name string, shorthand string, value []float64, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Float64SliceP(name, shorthand, value, usage)
	return b
}

// WithFloat64SliceVarPersistentFlag defines a float64Slice flag with specified
// name, default value, and usage string. The argument p points to a []float64
// variable in which to store the value of the flag.
func (b *CobraCmdBuilder) WithFloat64SliceVarPersistentFlag(variable *[]float64, name string, value []float64, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Float64SliceVar(variable, name, value, usage)
	return b
}

// WithFloat64SliceVarPPersistentFlag is like Float64SliceVar, but accepts a
// shorthand letter that can be used after a single dash.
func (b *CobraCmdBuilder) WithFloat64SliceVarPPersistentFlag(variable *[]float64, name string, shorthand string, value []float64, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Float64SliceVarP(variable, name, shorthand, value, usage)
	return b
}

// WithIntPersistentFlag defines an int flag with specified name, default value,
// and usage string. The return value is the address of an int variable that
// stores the value of the flag.
func (b *CobraCmdBuilder) WithIntPersistentFlag(name string, value int, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Int(name, value, usage)
	return b
}

// WithIntPPersistentFlag is like Int, but accepts a shorthand letter that can
// be used after a single dash.
func (b *CobraCmdBuilder) WithIntPPersistentFlag(name string, shorthand string, value int, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().IntP(name, shorthand, value, usage)
	return b
}

// WithIntVarPersistentFlag defines an int flag with specified name, default
// value, and usage string. The argument p points to an int variable in which to
// store the value of the flag.
func (b *CobraCmdBuilder) WithIntVarPersistentFlag(variable *int, name string, value int, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().IntVar(variable, name, value, usage)
	return b
}

// WithIntVarPPersistentFlag is like IntVar, but accepts a shorthand letter that
// can be used after a single dash.
func (b *CobraCmdBuilder) WithIntVarPPersistentFlag(variable *int, name string, shorthand string, value int, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().IntVarP(variable, name, shorthand, value, usage)
	return b
}

// WithIntSlicePersistentFlag defines a []int flag with specified name, default
// value, and usage string. The return value is the address of a []int variable
// that stores the value of the flag.
func (b *CobraCmdBuilder) WithIntSlicePersistentFlag(name string, value []int, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().IntSlice(name, value, usage)
	return b
}

// WithIntSlicePPersistentFlag is like IntSlice, but accepts a shorthand letter
// that can be used after a single dash.
func (b *CobraCmdBuilder) WithIntSlicePPersistentFlag(name string, shorthand string, value []int, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().IntSliceP(name, shorthand, value, usage)
	return b
}

// WithIntSliceVarPersistentFlag defines a intSlice flag with specified name,
// default value, and usage string. The argument p points to a []int variable in
// which to store the value of the flag.
func (b *CobraCmdBuilder) WithIntSliceVarPersistentFlag(variable *[]int, name string, value []int, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().IntSliceVar(variable, name, value, usage)
	return b
}

// WithIntSliceVarPPersistentFlag is like IntSliceVar, but accepts a shorthand
// letter that can be used after a single dash.
func (b *CobraCmdBuilder) WithIntSliceVarPPersistentFlag(variable *[]int, name string, shorthand string, value []int, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().IntSliceVarP(variable, name, shorthand, value, usage)
	return b
}

// WithInt8PersistentFlag defines an int8 flag with specified name, default
// value, and usage string. The return value is the address of an int8 variable
// that stores the value of the flag.
func (b *CobraCmdBuilder) WithInt8PersistentFlag(name string, value int8, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Int8(name, value, usage)
	return b
}

// WithInt8PPersistentFlag is like Int8, but accepts a shorthand letter that can
// be used after a single dash.
func (b *CobraCmdBuilder) WithInt8PPersistentFlag(name string, shorthand string, value int8, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Int8P(name, shorthand, value, usage)
	return b
}

// WithInt8VarPersistentFlag defines an int8 flag with specified name, default
// value, and usage string. The argument p points to an int8 variable in which
// to store the value of the flag.
func (b *CobraCmdBuilder) WithInt8VarPersistentFlag(variable *int8, name string, value int8, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Int8Var(variable, name, value, usage)
	return b
}

// WithInt8VarPPersistentFlag is like Int8Var, but accepts a shorthand letter
// that can be used after a single dash.
func (b *CobraCmdBuilder) WithInt8VarPPersistentFlag(variable *int8, name string, shorthand string, value int8, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Int8VarP(variable, name, shorthand, value, usage)
	return b
}

// WithInt16PersistentFlag defines an int16 flag with specified name, default
// value, and usage string. The return value is the address of an int16 variable
// that stores the value of the flag.
func (b *CobraCmdBuilder) WithInt16PersistentFlag(name string, value int16, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Int16(name, value, usage)
	return b
}

// WithInt16PPersistentFlag is like Int16, but accepts a shorthand letter that
// can be used after a single dash.
func (b *CobraCmdBuilder) WithInt16PPersistentFlag(name string, shorthand string, value int16, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Int16P(name, shorthand, value, usage)
	return b
}

// WithInt16VarPersistentFlag defines an int16 flag with specified name, default
// value, and usage string. The argument p points to an int16 variable in which
// to store the value of the flag.
func (b *CobraCmdBuilder) WithInt16VarPersistentFlag(variable *int16, name string, value int16, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Int16Var(variable, name, value, usage)
	return b
}

// WithInt16VarPPersistentFlag is like Int16Var, but accepts a shorthand letter
// that can be used after a single dash.
func (b *CobraCmdBuilder) WithInt16VarPPersistentFlag(variable *int16, name string, shorthand string, value int16, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Int16VarP(variable, name, shorthand, value, usage)
	return b
}

// WithInt32PersistentFlag defines an int32 flag with specified name, default
// value, and usage string. The return value is the address of an int32 variable
// that stores the value of the flag.
func (b *CobraCmdBuilder) WithInt32PersistentFlag(name string, value int32, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Int32(name, value, usage)
	return b
}

// WithInt32PPersistentFlag is like Int32, but accepts a shorthand letter that
// can be used after a single dash.
func (b *CobraCmdBuilder) WithInt32PPersistentFlag(name string, shorthand string, value int32, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Int32P(name, shorthand, value, usage)
	return b
}

// WithInt32VarPersistentFlag defines an int32 flag with specified name, default
// value, and usage string. The argument p points to an int32 variable in which
// to store the value of the flag.
func (b *CobraCmdBuilder) WithInt32VarPersistentFlag(variable *int32, name string, value int32, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Int32Var(variable, name, value, usage)
	return b
}

// WithInt32VarPPersistentFlag is like Int32Var, but accepts a shorthand letter
// that can be used after a single dash.
func (b *CobraCmdBuilder) WithInt32VarPPersistentFlag(variable *int32, name string, shorthand string, value int32, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Int32VarP(variable, name, shorthand, value, usage)
	return b
}

// WithInt32SlicePersistentFlag defines a []int32 flag with specified name,
// default value, and usage string. The return value is the address of a []int32
// variable that stores the value of the flag.
func (b *CobraCmdBuilder) WithInt32SlicePersistentFlag(name string, value []int32, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Int32Slice(name, value, usage)
	return b
}

// WithInt32SlicePPersistentFlag is like Int32Slice, but accepts a shorthand
// letter that can be used after a single dash.
func (b *CobraCmdBuilder) WithInt32SlicePPersistentFlag(name string, shorthand string, value []int32, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Int32SliceP(name, shorthand, value, usage)
	return b
}

// WithInt32SliceVarPersistentFlag defines a int32Slice flag with specified
// name, default value, and usage string. The argument p points to a []int32
// variable in which to store the value of the flag.
func (b *CobraCmdBuilder) WithInt32SliceVarPersistentFlag(variable *[]int32, name string, value []int32, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Int32SliceVar(variable, name, value, usage)
	return b
}

// WithInt32SliceVarPPersistentFlag is like Int32SliceVar, but accepts a
// shorthand letter that can be used after a single dash.
func (b *CobraCmdBuilder) WithInt32SliceVarPPersistentFlag(variable *[]int32, name string, shorthand string, value []int32, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Int32SliceVarP(variable, name, shorthand, value, usage)
	return b
}

// WithInt64PersistentFlag defines an int64 flag with specified name, default
// value, and usage string. The return value is the address of an int64 variable
// that stores the value of the flag.
func (b *CobraCmdBuilder) WithInt64PersistentFlag(name string, value int64, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Int64(name, value, usage)
	return b
}

// WithInt64PPersistentFlag is like Int64, but accepts a shorthand letter that
// can be used after a single dash.
func (b *CobraCmdBuilder) WithInt64PPersistentFlag(name string, shorthand string, value int64, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Int64P(name, shorthand, value, usage)
	return b
}

// WithInt64VarPersistentFlag defines an int64 flag with specified name, default
// value, and usage string. The argument p points to an int64 variable in which
// to store the value of the flag.
func (b *CobraCmdBuilder) WithInt64VarPersistentFlag(variable *int64, name string, value int64, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Int64Var(variable, name, value, usage)
	return b
}

// WithInt64VarPPersistentFlag is like Int64Var, but accepts a shorthand letter
// that can be used after a single dash.
func (b *CobraCmdBuilder) WithInt64VarPPersistentFlag(variable *int64, name string, shorthand string, value int64, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Int64VarP(variable, name, shorthand, value, usage)
	return b
}

// WithInt64SlicePersistentFlag defines a []int64 flag with specified name,
// default value, and usage string. The return value is the address of a []int64
// variable that stores the value of the flag.
func (b *CobraCmdBuilder) WithInt64SlicePersistentFlag(name string, value []int64, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Int64Slice(name, value, usage)
	return b
}

// WithInt64SlicePPersistentFlag is like Int64Slice, but accepts a shorthand
// letter that can be used after a single dash.
func (b *CobraCmdBuilder) WithInt64SlicePPersistentFlag(name string, shorthand string, value []int64, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Int64SliceP(name, shorthand, value, usage)
	return b
}

// WithInt64SliceVarPersistentFlag defines a int64Slice flag with specified
// name, default value, and usage string. The argument p points to a []int64
// variable in which to store the value of the flag.
func (b *CobraCmdBuilder) WithInt64SliceVarPersistentFlag(variable *[]int64, name string, value []int64, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Int64SliceVar(variable, name, value, usage)
	return b
}

// WithInt64SliceVarPPersistentFlag is like Int64SliceVar, but accepts a
// shorthand letter that can be used after a single dash.
func (b *CobraCmdBuilder) WithInt64SliceVarPPersistentFlag(variable *[]int64, name string, shorthand string, value []int64, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Int64SliceVarP(variable, name, shorthand, value, usage)
	return b
}

// WithUintPersistentFlag defines a uint flag with specified name, default
// value, and usage string. The return value is the address of a uint variable
// that stores the value of the flag.
func (b *CobraCmdBuilder) WithUintPersistentFlag(name string, value uint, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Uint(name, value, usage)
	return b
}

// WithIntPPersistentFlag is like Uint, but accepts a shorthand letter that can
// be used after a single dash.
func (b *CobraCmdBuilder) WithUintPPersistentFlag(name string, shorthand string, value uint, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().UintP(name, shorthand, value, usage)
	return b
}

// WithUintVarPersistentFlag defines a uint flag with specified name, default
// value, and usage string. The argument p points to a uint variable in which to
// store the value of the flag.
func (b *CobraCmdBuilder) WithUintVarPersistentFlag(variable *uint, name string, value uint, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().UintVar(variable, name, value, usage)
	return b
}

// WithUintVarPPersistentFlag is like UintVar, but accepts a shorthand letter
// that can be used after a single dash.
func (b *CobraCmdBuilder) WithUintVarPPersistentFlag(variable *uint, name string, shorthand string, value uint, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().UintVarP(variable, name, shorthand, value, usage)
	return b
}

// WithUintSlicePersistentFlag defines a []uint flag with specified name,
// default value, and usage string. The return value is the address of a []uint
// variable that stores the value of the flag.
func (b *CobraCmdBuilder) WithUintSlicePersistentFlag(name string, value []uint, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().UintSlice(name, value, usage)
	return b
}

// WithUintSlicePPersistentFlag is like UintSlice, but accepts a shorthand
// letter that can be used after a single dash.
func (b *CobraCmdBuilder) WithUintSlicePPersistentFlag(name string, shorthand string, value []uint, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().UintSliceP(name, shorthand, value, usage)
	return b
}

// WithUintSliceVarPersistentFlag defines a uintSlice flag with specified name,
// default value, and usage string. The argument p points to a []uint variable
// in which to store the value of the flag.
func (b *CobraCmdBuilder) WithUintSliceVarPersistentFlag(variable *[]uint, name string, value []uint, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().UintSliceVar(variable, name, value, usage)
	return b
}

// WithUintSliceVarPPersistentFlag is like UintSliceVar, but accepts a shorthand
// letter that can be used after a single dash.
func (b *CobraCmdBuilder) WithUintSliceVarPPersistentFlag(variable *[]uint, name string, shorthand string, value []uint, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().UintSliceVarP(variable, name, shorthand, value, usage)
	return b
}

// WithUint8PersistentFlag defines a uint8 flag with specified name, default
// value, and usage string. The return value is the address of a uint8 variable
// that stores the value of the flag.
func (b *CobraCmdBuilder) WithUint8PersistentFlag(name string, value uint8, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Uint8(name, value, usage)
	return b
}

// WithUint8PPersistentFlag is like Uint8, but accepts a shorthand letter that
// can be used after a single dash.
func (b *CobraCmdBuilder) WithUint8PPersistentFlag(name string, shorthand string, value uint8, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Uint8P(name, shorthand, value, usage)
	return b
}

// WithUint8VarPersistentFlag defines a uint8 flag with specified name, default
// value, and usage string. The argument p points to a uint8 variable in which
// to store the value of the flag.
func (b *CobraCmdBuilder) WithUint8VarPersistentFlag(variable *uint8, name string, value uint8, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Uint8Var(variable, name, value, usage)
	return b
}

// WithUint8VarPPersistentFlag is like Uint8Var, but accepts a shorthand letter
// that can be used after a single dash.
func (b *CobraCmdBuilder) WithUint8VarPPersistentFlag(variable *uint8, name string, shorthand string, value uint8, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Uint8VarP(variable, name, shorthand, value, usage)
	return b
}

// WithUint16PersistentFlag defines a uint flag with specified name, default
// value, and usage string. The return value is the address of a uint variable
// that stores the value of the flag.
func (b *CobraCmdBuilder) WithUint16PersistentFlag(name string, value uint16, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Uint16(name, value, usage)
	return b
}

// WithUint16PPersistentFlag is like Uint16, but accepts a shorthand letter that
// can be used after a single dash.
func (b *CobraCmdBuilder) WithUint16PPersistentFlag(name string, shorthand string, value uint16, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Uint16P(name, shorthand, value, usage)
	return b
}

// WithUint16VarPersistentFlag defines a uint flag with specified name, default
// value, and usage string. The argument p points to a uint variable in which to
// store the value of the flag.
func (b *CobraCmdBuilder) WithUint16VarPersistentFlag(variable *uint16, name string, value uint16, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Uint16Var(variable, name, value, usage)
	return b
}

// WithUint16VarPPersistentFlag is like Uint16Var, but accepts a shorthand
// letter that can be used after a single dash.
func (b *CobraCmdBuilder) WithUint16VarPPersistentFlag(variable *uint16, name string, shorthand string, value uint16, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Uint16VarP(variable, name, shorthand, value, usage)
	return b
}

// WithUint32PersistentFlag defines a uint32 flag with specified name, default
// value, and usage string. The return value is the address of a uint32 variable
// that stores the value of the flag.
func (b *CobraCmdBuilder) WithUint32PersistentFlag(name string, value uint32, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Uint32(name, value, usage)
	return b
}

// WithUint32PPersistentFlag is like Uint32, but accepts a shorthand letter that
// can be used after a single dash.
func (b *CobraCmdBuilder) WithUint32PPersistentFlag(name string, shorthand string, value uint32, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Uint32P(name, shorthand, value, usage)
	return b
}

// WithUint32VarPersistentFlag defines a uint32 flag with specified name,
// default value, and usage string. The argument p points to a uint32 variable
// in which to store the value of the flag.
func (b *CobraCmdBuilder) WithUint32VarPersistentFlag(variable *uint32, name string, value uint32, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Uint32Var(variable, name, value, usage)
	return b
}

// WithUint32VarPPersistentFlag is like Uint32Var, but accepts a shorthand
// letter that can be used after a single dash.
func (b *CobraCmdBuilder) WithUint32VarPPersistentFlag(variable *uint32, name string, shorthand string, value uint32, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Uint32VarP(variable, name, shorthand, value, usage)
	return b
}

// WithUint64PersistentFlag defines a uint64 flag with specified name, default
// value, and usage string. The return value is the address of a uint64 variable
// that stores the value of the flag.
func (b *CobraCmdBuilder) WithUint64PersistentFlag(name string, value uint64, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Uint64(name, value, usage)
	return b
}

// WithUint64PPersistentFlag is like Uint64, but accepts a shorthand letter that
// can be used after a single dash.
func (b *CobraCmdBuilder) WithUint64PPersistentFlag(name string, shorthand string, value uint64, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Uint64P(name, shorthand, value, usage)
	return b
}

// WithUint64VarPersistentFlag defines a uint64 flag with specified name,
// default value, and usage string. The argument p points to a uint64 variable
// in which to store the value of the flag.
func (b *CobraCmdBuilder) WithUint64VarPersistentFlag(variable *uint64, name string, value uint64, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Uint64Var(variable, name, value, usage)
	return b
}

// WithUint64VarPPersistentFlag is like Uint64Var, but accepts a shorthand
// letter that can be used after a single dash.
func (b *CobraCmdBuilder) WithUint64VarPPersistentFlag(variable *uint64, name string, shorthand string, value uint64, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Uint64VarP(variable, name, shorthand, value, usage)
	return b
}

// WithStringPersistentFlag defines a string flag with specified name, default
// value, and usage string. The return value is the address of a string variable
// that stores the value of the flag.
func (b *CobraCmdBuilder) WithStringPersistentFlag(name string, value string, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().String(name, value, usage)
	return b
}

// WithStringPPersistentFlag is like String, but accepts a shorthand letter that
// can be used after a single dash.
func (b *CobraCmdBuilder) WithStringPPersistentFlag(name string, shorthand string, value string, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().StringP(name, shorthand, value, usage)
	return b
}

// WithStringVarPersistentFlag defines a string flag with specified name,
// default value, and usage string. The argument p points to a string variable
// in which to store the value of the flag.
func (b *CobraCmdBuilder) WithStringVarPersistentFlag(variable *string, name string, value string, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().StringVar(variable, name, value, usage)
	return b
}

// WithStringVarPPersistentFlag is like StringVar, but accepts a shorthand
// letter that can be used after a single dash.
func (b *CobraCmdBuilder) WithStringVarPPersistentFlag(variable *string, name string, shorthand string, value string, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().StringVarP(variable, name, shorthand, value, usage)
	return b
}

// WithStringSlicePersistentFlag defines a string flag with specified name,
// default value, and usage string. The return value is the address of a
// []string variable that stores the value of the flag. Compared to StringArray
// flags, StringSlice flags take comma-separated value as arguments and split
// them accordingly.
// For example:
//
//	--ss="v1,v2" --ss="v3"
//
// will result in
//
//	[]string{"v1", "v2", "v3"}
func (b *CobraCmdBuilder) WithStringSlicePersistentFlag(name string, value []string, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().StringSlice(name, value, usage)
	return b
}

// WithStringSlicePPersistentFlag is like StringSlice, but accepts a shorthand
// letter that can be used after a single dash.
func (b *CobraCmdBuilder) WithStringSlicePPersistentFlag(name string, shorthand string, value []string, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().StringSliceP(name, shorthand, value, usage)
	return b
}

// WithStringSliceVarPersistentFlag defines a string flag with specified name,
// default value, and usage string. The argument p points to a []string variable
// in which to store the value of the flag. Compared to StringArray flags,
// StringSlice flags take comma-separated value as arguments and split them
// accordingly.
// For example:
//
//	--ss="v1,v2" --ss="v3"
//
// will result in
//
//	[]string{"v1", "v2", "v3"}
func (b *CobraCmdBuilder) WithStringSliceVarPersistentFlag(variable *[]string, name string, value []string, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().StringSliceVar(variable, name, value, usage)
	return b
}

// WithStringSliceVarPPersistentFlag is like StringSliceVar, but accepts a
// shorthand letter that can be used after a single dash.
func (b *CobraCmdBuilder) WithStringSliceVarPPersistentFlag(variable *[]string, name string, shorthand string, value []string, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().StringSliceVarP(variable, name, shorthand, value, usage)
	return b
}

// WithStringArrayPersistentFlag defines a string flag with specified name,
// default value, and usage string. The return value is the address of a
// []string variable that stores the value of the flag. The value of each
// argument will not try to be separated by comma. Use a StringSlice for that.
func (b *CobraCmdBuilder) WithStringArrayPersistentFlag(name string, value []string, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().StringArray(name, value, usage)
	return b
}

// WithStringArrayPPersistentFlag is like StringArray, but accepts a shorthand
// letter that can be used after a single dash.
func (b *CobraCmdBuilder) WithStringArrayPPersistentFlag(name string, shorthand string, value []string, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().StringArrayP(name, shorthand, value, usage)
	return b
}

// WithStringArrayVarPersistentFlag defines a string flag with specified name,
// default value, and usage string. The argument p points to a []string variable
// in which to store the values of the multiple flags. The value of each
// argument will not try to be separated by comma. Use a StringSlice for that.
func (b *CobraCmdBuilder) WithStringArrayVarPersistentFlag(variable *[]string, name string, value []string, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().StringArrayVar(variable, name, value, usage)
	return b
}

// WithStringArrayVarPPersistentFlag is like StringArrayVar, but accepts a
// shorthand letter that can be used after a single dash.
func (b *CobraCmdBuilder) WithStringArrayVarPPersistentFlag(variable *[]string, name string, shorthand string, value []string, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().StringArrayVarP(variable, name, shorthand, value, usage)
	return b
}

// WithStringToIntPersistentFlag defines a string flag with specified name,
// default value, and usage string. The return value is the address of a
// map[string]int variable that stores the value of the flag. The value of each
// argument will not try to be separated by comma
func (b *CobraCmdBuilder) WithStringToIntPersistentFlag(name string, value map[string]int, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().StringToInt(name, value, usage)
	return b
}

// WithStringToIntPPersistentFlag is like StringToInt, but accepts a shorthand
// letter that can be used after a single dash.
func (b *CobraCmdBuilder) WithStringToIntPPersistentFlag(name string, shorthand string, value map[string]int, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().StringToIntP(name, shorthand, value, usage)
	return b
}

// WithStringToIntVarPersistentFlag defines a string flag with specified name,
// default value, and usage string. The argument p points to a map[string]int
// variable in which to store the values of the multiple flags. The value of
// each argument will not try to be separated by comma
func (b *CobraCmdBuilder) WithStringToIntVarPersistentFlag(variable *map[string]int, name string, value map[string]int, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().StringToIntVar(variable, name, value, usage)
	return b
}

// WithStringToIntVarPPersistentFlag is like StringToIntVar, but accepts a
// shorthand letter that can be used after a single dash.
func (b *CobraCmdBuilder) WithStringToIntVarPPersistentFlag(variable *map[string]int, name string, shorthand string, value map[string]int, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().StringToIntVarP(variable, name, shorthand, value, usage)
	return b
}

// WithStringToInt64PersistentFlag defines a string flag with specified name,
// default value, and usage string. The return value is the address of a
// map[string]int64 variable that stores the value of the flag. The value of
// each argument will not try to be separated by comma
func (b *CobraCmdBuilder) WithStringToInt64PersistentFlag(name string, value map[string]int64, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().StringToInt64(name, value, usage)
	return b
}

// WithStringToInt64PPersistentFlag is like StringToInt64, but accepts a
// shorthand letter that can be used after a single dash.
func (b *CobraCmdBuilder) WithStringToInt64PPersistentFlag(name string, shorthand string, value map[string]int64, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().StringToInt64P(name, shorthand, value, usage)
	return b
}

// WithStringToInt64VarPersistentFlag defines a string flag with specified name,
// default value, and usage string. The argument p point64s to a
// map[string]int64 variable in which to store the values of the multiple flags.
// The value of each argument will not try to be separated by comma
func (b *CobraCmdBuilder) WithStringToInt64VarPersistentFlag(variable *map[string]int64, name string, value map[string]int64, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().StringToInt64Var(variable, name, value, usage)
	return b
}

// WithStringToInt64VarPPersistentFlag is like StringToInt64Var, but accepts a
// shorthand letter that can be used after a single dash.
func (b *CobraCmdBuilder) WithStringToInt64VarPPersistentFlag(variable *map[string]int64, name string, shorthand string, value map[string]int64, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().StringToInt64VarP(variable, name, shorthand, value, usage)
	return b
}

// WithIPPersistentFlag defines an net.IP flag with specified name, default
// value, and usage string. The return value is the address of an net.IP
// variable that stores the value of the flag.
func (b *CobraCmdBuilder) WithIPPersistentFlag(name string, value net.IP, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().IP(name, value, usage)
	return b
}

// WithIPPPersistentFlag is like IP, but accepts a shorthand letter that can be
// used after a single dash.
func (b *CobraCmdBuilder) WithIPPPersistentFlag(name string, shorthand string, value net.IP, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().IPP(name, shorthand, value, usage)
	return b
}

// WithIPVarPersistentFlag defines an net.IP flag with specified name, default
// value, and usage string. The argument p points to an net.IP variable in which
// to store the value of the flag.
func (b *CobraCmdBuilder) WithIPVarPersistentFlag(variable *net.IP, name string, value net.IP, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().IPVar(variable, name, value, usage)
	return b
}

// WithIPVarPPersistentFlag is like IPVar, but accepts a shorthand letter that
// can be used after a single dash.
func (b *CobraCmdBuilder) WithIPVarPPersistentFlag(variable *net.IP, name string, shorthand string, value net.IP, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().IPVarP(variable, name, shorthand, value, usage)
	return b
}

// WithIPSlicePersistentFlag defines a []net.IP flag with specified name,
// default value, and usage string. The return value is the address of a
// []net.IP variable that stores the value of that flag.
func (b *CobraCmdBuilder) WithIPSlicePersistentFlag(name string, value []net.IP, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().IPSlice(name, value, usage)
	return b
}

// WithIPSlicePPersistentFlag is like IPSlice, but accepts a shorthand letter
// that can be used after a single dash.
func (b *CobraCmdBuilder) WithIPSlicePPersistentFlag(name string, shorthand string, value []net.IP, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().IPSliceP(name, shorthand, value, usage)
	return b
}

// WithIPSliceVarPersistentFlag defines a ipSlice flag with specified name,
// default value, and usage string. The argument p points to a []net.IP variable
// in which to store the value of the flag.
func (b *CobraCmdBuilder) WithIPSliceVarPersistentFlag(variable *[]net.IP, name string, value []net.IP, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().IPSliceVar(variable, name, value, usage)
	return b
}

// WithIPSliceVarPPersistentFlag is like IPSliceVar, but accepts a shorthand
// letter that can be used after a single dash.
func (b *CobraCmdBuilder) WithIPSliceVarPPersistentFlag(variable *[]net.IP, name string, shorthand string, value []net.IP, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().IPSliceVarP(variable, name, shorthand, value, usage)
	return b
}

// WithIPMaskPersistentFlag defines an net.IPMask flag with specified name,
// default value, and usage string. The return value is the address of an
// net.IPMask variable that stores the value of the flag.
func (b *CobraCmdBuilder) WithIPMaskPersistentFlag(name string, value net.IPMask, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().IPMask(name, value, usage)
	return b
}

// WithIPMaskPPersistentFlag is like IPMask, but accepts a shorthand letter that
// can be used after a single dash.
func (b *CobraCmdBuilder) WithIPMaskPPersistentFlag(name string, shorthand string, value net.IPMask, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().IPMaskP(name, shorthand, value, usage)
	return b
}

// WithIPMaskVarPersistentFlag defines an net.IPMask flag with specified name,
// default value, and usage string. The argument p points to an net.IPMask
// variable in which to store the value of the flag.
func (b *CobraCmdBuilder) WithIPMaskVarPersistentFlag(variable *net.IPMask, name string, value net.IPMask, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().IPMaskVar(variable, name, value, usage)
	return b
}

// WithIPMaskVarPPersistentFlag is like IPMaskVar, but accepts a shorthand
// letter that can be used after a single dash.
func (b *CobraCmdBuilder) WithIPMaskVarPPersistentFlag(variable *net.IPMask, name string, shorthand string, value net.IPMask, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().IPMaskVarP(variable, name, shorthand, value, usage)
	return b
}

// WithIPNetPersistentFlag defines an net.IPNet flag with specified name,
// default value, and usage string. The return value is the address of an
// net.IPNet variable that stores the value of the flag.
func (b *CobraCmdBuilder) WithIPNetPersistentFlag(name string, value net.IPNet, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().IPNet(name, value, usage)
	return b
}

// WithIPNetPPersistentFlag is like IPNet, but accepts a shorthand letter that
// can be used after a single dash.
func (b *CobraCmdBuilder) WithIPNetPPersistentFlag(name string, shorthand string, value net.IPNet, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().IPNetP(name, shorthand, value, usage)
	return b
}

// WithIPNetVarPersistentFlag defines an net.IPNet flag with specified name,
// default value, and usage string. The argument p points to an net.IPNet
// variable in which to store the value of the flag.
func (b *CobraCmdBuilder) WithIPNetVarPersistentFlag(variable *net.IPNet, name string, value net.IPNet, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().IPNetVar(variable, name, value, usage)
	return b
}

// WithIPNetVarPPersistentFlag is like IPNetVar, but accepts a shorthand letter
// that can be used after a single dash.
func (b *CobraCmdBuilder) WithIPNetVarPPersistentFlag(variable *net.IPNet, name string, shorthand string, value net.IPNet, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().IPNetVarP(variable, name, shorthand, value, usage)
	return b
}

// WithVarPersistentFlag defines a flag with the specified name and usage
// string. The type and value of the flag are represented by the first argument,
// of type Value, which typically holds a user-defined implementation of Value.
// For instance, the caller could create a flag that turns a comma-separated
// string into a slice of strings by giving the slice the methods of Value; in
// particular, Set would decompose the comma-separated string into the slice.
func (b *CobraCmdBuilder) WithVarPersistentFlag(value pflag.Value, name string, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().Var(value, name, usage)
	return b
}

// WithVarPPersistentFlag is like Var, but accepts a shorthand letter that can
// be used after a single dash.
func (b *CobraCmdBuilder) WithVarPPersistentFlag(value pflag.Value, name string, shorthand string, usage string) *CobraCmdBuilder {
	b.cmd.PersistentFlags().VarP(value, name, shorthand, usage)
	return b
}

// MarkFlagHidden sets a flag to 'hidden' in your program. It will continue to
// function but will not show up in help or usage messages.
func (b *CobraCmdBuilder) MarkPersistentFlagHidden(name string) *CobraCmdBuilder {
	err := b.cmd.PersistentFlags().MarkHidden(name)
	if err != nil {
		panic(err)
	}
	return b
}

// MarkFlagDeprecated indicated that a flag is deprecated in your program. It
// will continue to function but will not show up in help or usage messages.
// Using this flag will also print the given usageMessage.
func (b *CobraCmdBuilder) MarkPersistentFlagDeprecated(name string, usage string) *CobraCmdBuilder {
	err := b.cmd.PersistentFlags().MarkDeprecated(name, usage)
	if err != nil {
		panic(err)
	}
	return b
}

// MarkFlagShorthandDeprecated will mark the shorthand of a flag deprecated in
// your program. It will continue to function but will not show up in help or
// usage messages. Using this flag will also print the given usageMessage.
func (b *CobraCmdBuilder) MarkPersistentFlagShorthandDeprecated(name string, usage string) *CobraCmdBuilder {
	err := b.cmd.PersistentFlags().MarkShorthandDeprecated(name, usage)
	if err != nil {
		panic(err)
	}
	return b
}

// WithFlagSet adds one FlagSet to another. If a flag is already present in f
// the flag from newSet will be ignored.
func (b *CobraCmdBuilder) WithFlagSet(flagset *pflag.FlagSet) *CobraCmdBuilder {
	b.cmd.Flags().AddFlagSet(flagset)
	return b
}

// WithPersistentFlagSet adds one FlagSet to another. If a flag is already
// present in f the flag from newSet will be ignored.
func (b *CobraCmdBuilder) WithPersistentFlagSet(flagset *pflag.FlagSet) *CobraCmdBuilder {
	b.cmd.PersistentFlags().AddFlagSet(flagset)
	return b
}

// ToBoaCmdBuilder returns a BoaCmdBuilder from a CobraCmdBuilder
func (b *CobraCmdBuilder) ToBoaCmdBuilder() *BoaCmdBuilder {
	return &BoaCmdBuilder{
		b,
		&Command{
			b.cmd,
			[]Option{},
			[]Profile{},
		},
	}
}

// BuildBoaCmd returns a boa Command from a CobraCmdBuilder
func (b *CobraCmdBuilder) BuildBoaCmd() *Command {
	return &Command{
		b.cmd,
		[]Option{},
		[]Profile{},
	}
}

// Build returns a cobra.Command from a CobraCmdBuilder
func (b *CobraCmdBuilder) Build() *cobra.Command {
	return b.cmd
}
