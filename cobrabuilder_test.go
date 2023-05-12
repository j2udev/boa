package boa

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCobraCmd struct {
	mock.Mock
	*cobra.Command
}

func TestCobraCmdBuilder(t *testing.T) {
	// common values to use for struct literal cobra.Command vs CobraCmdBuilder
	// cobra.Command
	aliases := []string{"alias1", "alias2"}
	suggestFor := []string{"cmd1", "cmd2"}
	short := "short desc"
	long := "long desc"
	groupId := "group ID"
	example := `example 1`
	validArgs := []string{"arg1", "arg2"}
	validArgsFunc := func(*cobra.Command, []string, string) ([]string, cobra.ShellCompDirective) {
		return []string{"test"}, cobra.ShellCompDirectiveDefault
	}
	args := cobra.MatchAll(cobra.MaximumNArgs(1), cobra.OnlyValidArgs)
	argAliases := []string{"argAlias1", "argAlias2"}
	bashCompFunc := "compFunc"
	deprecationMsg := "cmd is deprecated"
	annotations := make(map[string]string)
	version := "0.1.0"
	runFunc := func(*cobra.Command, []string) {}
	runEFunc := func(*cobra.Command, []string) error {
		return fmt.Errorf("test")
	}
	errWhitelist := cobra.FParseErrWhitelist{}
	options := cobra.CompletionOptions{}
	traverseChildren := true
	hidden := true
	silenceErrors := true
	silenceUsage := true
	disableFlagParsing := true
	disableFlagsInUseLine := true
	disableAutoGenTag := true
	disableSuggestions := true
	suggestionMinDist := 2
	subCommands := []*cobra.Command{{Use: "a"}, {Use: "b"}}
	usageTemplate := "usage"
	helpTemplate := "help"
	usageFunc := func(*cobra.Command) error { return errors.New("") }
	helpFunc := func(*cobra.Command, []string) {}

	// common values to use for mock cobra.Command flagset vs CobraCmdBuilder flagset
	boolName := "boolFlag"
	boolVarName := "boolVarFlag"
	boolPName := "boolPFlag"
	boolVarPName := "boolVarPFlag"
	boolVar := false
	boolShorthand := "a"
	boolVarShorthand := "A"
	boolDefault := false
	boolUsage := "bool usage"
	boolSliceName := "boolSliceFlag"
	boolSliceVarName := "boolSliceVarFlag"
	boolSlicePName := "boolSlicePFlag"
	boolSliceVarPName := "boolSliceVarPFlag"
	boolSliceVar := []bool{}
	boolSliceShorthand := "b"
	boolSliceVarShorthand := "B"
	boolSliceDefault := []bool{false}
	boolSliceUsage := "bool slice usage"
	stringName := "stringFlag"
	stringVarName := "stringVarFlag"
	stringPName := "stringPFlag"
	stringVarPName := "stringVarPFlag"
	stringSliceName := "stringSliceFlag"
	stringSliceVarName := "stringSliceVarFlag"
	stringSlicePName := "stringSlicePFlag"
	stringSliceVarPName := "stringSliceVarPFlag"
	stringVar := "t"
	stringShorthand := "c"
	stringVarShorthand := "C"
	stringDefault := "test"
	stringUsage := "string usage"
	stringSliceVar := []string{"t"}
	stringSliceShorthand := "d"
	stringSliceVarShorthand := "D"
	stringSliceDefault := []string{"test"}
	stringSliceUsage := "string slice usage"

	// define mock cobra.Command using struct literal to compare to the resulting
	// cobra.Command from the CobraCmdBuilder
	mockCmd := &MockCobraCmd{
		Command: &cobra.Command{
			Aliases:                    aliases,
			SuggestFor:                 suggestFor,
			Short:                      short,
			Long:                       long,
			GroupID:                    groupId,
			Example:                    example,
			ValidArgs:                  validArgs,
			ValidArgsFunction:          validArgsFunc,
			Args:                       args,
			ArgAliases:                 argAliases,
			BashCompletionFunction:     bashCompFunc,
			Deprecated:                 deprecationMsg,
			Annotations:                annotations,
			Version:                    version,
			PersistentPreRun:           runFunc,
			PersistentPreRunE:          runEFunc,
			PreRun:                     runFunc,
			PreRunE:                    runEFunc,
			Run:                        runFunc,
			RunE:                       runEFunc,
			PostRun:                    runFunc,
			PostRunE:                   runEFunc,
			PersistentPostRun:          runFunc,
			PersistentPostRunE:         runEFunc,
			FParseErrWhitelist:         errWhitelist,
			CompletionOptions:          options,
			TraverseChildren:           traverseChildren,
			Hidden:                     hidden,
			SilenceErrors:              silenceErrors,
			SilenceUsage:               silenceUsage,
			DisableFlagParsing:         disableFlagParsing,
			DisableFlagsInUseLine:      disableFlagsInUseLine,
			DisableAutoGenTag:          disableAutoGenTag,
			DisableSuggestions:         disableSuggestions,
			SuggestionsMinimumDistance: suggestionMinDist,
		},
	}
	// cover the additional cobra.Command methods that the CobraCmdBuilder wraps
	mockCmd.AddCommand(subCommands...)
	mockCmd.SetUsageTemplate(usageTemplate)
	mockCmd.SetHelpTemplate(helpTemplate)
	mockCmd.SetUsageFunc(usageFunc)
	mockCmd.SetHelpFunc(helpFunc)

	// define a mock pflag.Flagset to compare to the cobra.Command.Flags() flagset
	// that results from using the builder methods
	mockFs := pflag.NewFlagSet("mockfs", pflag.ContinueOnError)
	mockFs.Bool(boolName, boolDefault, boolUsage)
	mockFs.BoolVar(&boolVar, boolVarName, boolDefault, boolUsage)
	mockFs.BoolP(boolPName, boolShorthand, boolDefault, boolUsage)
	mockFs.BoolVarP(&boolVar, boolVarPName, boolVarShorthand, boolDefault, boolUsage)
	mockFs.BoolSlice(boolSliceName, boolSliceDefault, boolSliceUsage)
	mockFs.BoolSliceVar(&boolSliceVar, boolSliceVarName, boolSliceDefault, boolSliceUsage)
	mockFs.BoolSliceP(boolSlicePName, boolSliceShorthand, boolSliceDefault, boolSliceUsage)
	mockFs.BoolSliceVarP(&boolSliceVar, boolSliceVarPName, boolSliceVarShorthand, boolSliceDefault, boolSliceUsage)
	mockFs.String(stringName, stringDefault, stringUsage)
	mockFs.StringVar(&stringVar, stringVarName, stringDefault, stringUsage)
	mockFs.StringP(stringPName, stringShorthand, stringDefault, stringUsage)
	mockFs.StringVarP(&stringVar, stringVarPName, stringVarShorthand, stringDefault, stringUsage)
	mockFs.StringSlice(stringSliceName, stringSliceDefault, stringSliceUsage)
	mockFs.StringSliceVar(&stringSliceVar, stringSliceVarName, stringSliceDefault, stringSliceUsage)
	mockFs.StringSliceP(stringSlicePName, stringSliceShorthand, stringSliceDefault, stringSliceUsage)
	mockFs.StringSliceVarP(&stringSliceVar, stringSliceVarPName, stringSliceVarShorthand, stringSliceDefault, stringSliceUsage)
	mockCmd.Flags().AddFlagSet(mockFs)

	// Create a new cobra.Command using the builder
	cmd := NewCobraCmd("test").
		WithAliases(aliases).
		SuggestFor(suggestFor).
		WithShortDescription(short).
		WithLongDescription(long).
		WithGroupID(groupId).
		WithExample(example).
		WithValidArgs(validArgs).
		WithValidArgsFunction(validArgsFunc).
		WithArgs(args).
		WithArgAliases(argAliases).
		WithBashCompletionFunction(bashCompFunc).
		Deprecated(deprecationMsg).
		WithAnnotations(annotations).
		WithVersion(version).
		WithPersistentPreRunFunc(runFunc).
		WithPersistentPreRunEFunc(runEFunc).
		WithPreRunFunc(runFunc).
		WithPreRunEFunc(runEFunc).
		WithRunFunc(runFunc).
		WithRunEFunc(runEFunc).
		WithPostRunFunc(runFunc).
		WithPostRunEFunc(runEFunc).
		WithPersistentPostRunFunc(runFunc).
		WithPersistentPostRunEFunc(runEFunc).
		WithFParseErrWhitelist(errWhitelist).
		WithCompletionOptions(options).
		WithSuggestionsMinimumDistance(suggestionMinDist).
		WithSubCommands(subCommands...).
		WithUsageTemplate(usageTemplate).
		WithHelpTemplate(helpTemplate).
		WithUsageFunc(usageFunc).
		WithHelpFunc(helpFunc).
		TraverseChildren().
		Hidden().
		SilenceErrors().
		SilenceUsage().
		DisableFlagParsing().
		DisableFlagsInUseLine().
		DisableAutoGenTag().
		DisableSuggestions().
		WithBoolFlag(boolName, boolDefault, boolUsage).
		WithBoolVarFlag(&boolVar, boolVarName, boolDefault, boolUsage).
		WithBoolPFlag(boolPName, boolShorthand, boolDefault, boolUsage).
		WithBoolVarPFlag(&boolVar, boolVarPName, boolVarShorthand, boolDefault, boolUsage).
		WithBoolSliceFlag(boolSliceName, boolSliceDefault, boolSliceUsage).
		WithBoolSliceVarFlag(&boolSliceVar, boolSliceVarName, boolSliceDefault, boolSliceUsage).
		WithBoolSlicePFlag(boolSlicePName, boolSliceShorthand, boolSliceDefault, boolSliceUsage).
		WithBoolSliceVarPFlag(&boolSliceVar, boolSliceVarPName, boolSliceVarShorthand, boolSliceDefault, boolSliceUsage).
		WithStringFlag(stringName, stringDefault, stringUsage).
		WithStringVarFlag(&stringVar, stringVarName, stringDefault, stringUsage).
		WithStringPFlag(stringPName, stringShorthand, stringDefault, stringUsage).
		WithStringVarPFlag(&stringVar, stringVarPName, stringVarShorthand, stringDefault, stringUsage).
		WithStringSliceFlag(stringSliceName, stringSliceDefault, stringSliceUsage).
		WithStringSliceVarFlag(&stringSliceVar, stringSliceVarName, stringSliceDefault, stringSliceUsage).
		WithStringSlicePFlag(stringSlicePName, stringSliceShorthand, stringSliceDefault, stringSliceUsage).
		WithStringSliceVarPFlag(&stringSliceVar, stringSliceVarPName, stringSliceVarShorthand, stringSliceDefault, stringSliceUsage).
		Build()

	//Field Equality
	assert.Equal(t, mockCmd.Aliases, cmd.Aliases)
	assert.Equal(t, mockCmd.SuggestFor, cmd.SuggestFor)
	assert.Equal(t, mockCmd.Short, cmd.Short)
	assert.Equal(t, mockCmd.Long, cmd.Long)
	assert.Equal(t, mockCmd.GroupID, cmd.GroupID)
	assert.Equal(t, mockCmd.Example, cmd.Example)
	assert.Equal(t, mockCmd.ValidArgs, cmd.ValidArgs)
	assert.Equal(t, mockCmd.ArgAliases, cmd.ArgAliases)
	assert.Equal(t, mockCmd.BashCompletionFunction, cmd.BashCompletionFunction)
	assert.Equal(t, mockCmd.Deprecated, cmd.Deprecated)
	assert.Equal(t, mockCmd.Annotations, cmd.Annotations)
	assert.Equal(t, mockCmd.Version, cmd.Version)
	assert.Equal(t, mockCmd.FParseErrWhitelist, cmd.FParseErrWhitelist)
	assert.Equal(t, mockCmd.CompletionOptions, cmd.CompletionOptions)
	assert.Equal(t, mockCmd.TraverseChildren, cmd.TraverseChildren)
	assert.Equal(t, mockCmd.Hidden, cmd.Hidden)
	assert.Equal(t, mockCmd.SilenceErrors, cmd.SilenceErrors)
	assert.Equal(t, mockCmd.SilenceUsage, cmd.SilenceUsage)
	assert.Equal(t, mockCmd.DisableFlagParsing, cmd.DisableFlagParsing)
	assert.Equal(t, mockCmd.DisableFlagsInUseLine, cmd.DisableFlagsInUseLine)
	assert.Equal(t, mockCmd.DisableAutoGenTag, cmd.DisableAutoGenTag)
	assert.Equal(t, mockCmd.DisableSuggestions, cmd.DisableSuggestions)
	// Function Equality
	assert.Equal(t, getFuncName(mockCmd.ValidArgsFunction), getFuncName(cmd.ValidArgsFunction))
	assert.Equal(t, getFuncName(mockCmd.Args), getFuncName(cmd.Args))
	assert.Equal(t, getFuncName(mockCmd.PersistentPreRun), getFuncName(cmd.PersistentPreRun))
	assert.Equal(t, getFuncName(mockCmd.PersistentPreRunE), getFuncName(cmd.PersistentPreRunE))
	assert.Equal(t, getFuncName(mockCmd.PreRun), getFuncName(cmd.PreRun))
	assert.Equal(t, getFuncName(mockCmd.PreRunE), getFuncName(cmd.PreRunE))
	assert.Equal(t, getFuncName(mockCmd.Run), getFuncName(cmd.Run))
	assert.Equal(t, getFuncName(mockCmd.RunE), getFuncName(cmd.RunE))
	assert.Equal(t, getFuncName(mockCmd.PostRun), getFuncName(cmd.PostRun))
	assert.Equal(t, getFuncName(mockCmd.PostRunE), getFuncName(cmd.PostRunE))
	assert.Equal(t, getFuncName(mockCmd.PersistentPostRun), getFuncName(cmd.PersistentPostRun))
	assert.Equal(t, getFuncName(mockCmd.PersistentPostRunE), getFuncName(cmd.PersistentPostRunE))
	// Flag Equality
	// bool flags
	mockBoolFlag, _ := mockCmd.Flags().GetBool(boolName)
	cmdBoolFlag, _ := cmd.Flags().GetBool(boolName)
	assert.Equal(t, mockBoolFlag, cmdBoolFlag)
	mockBoolFlag, _ = mockCmd.Flags().GetBool(boolVarName)
	cmdBoolFlag, _ = cmd.Flags().GetBool(boolVarName)
	assert.Equal(t, mockBoolFlag, cmdBoolFlag)
	mockBoolFlag, _ = mockCmd.Flags().GetBool(boolPName)
	cmdBoolFlag, _ = cmd.Flags().GetBool(boolPName)
	assert.Equal(t, mockBoolFlag, cmdBoolFlag)
	mockBoolFlag, _ = mockCmd.Flags().GetBool(boolVarPName)
	cmdBoolFlag, _ = cmd.Flags().GetBool(boolVarPName)
	assert.Equal(t, mockBoolFlag, cmdBoolFlag)
	// bool slice flags
	mockBoolSliceFlag, _ := mockCmd.Flags().GetBoolSlice(boolSliceName)
	cmdBoolSliceFlag, _ := cmd.Flags().GetBoolSlice(boolSliceName)
	assert.Equal(t, mockBoolSliceFlag, cmdBoolSliceFlag)
	mockBoolSliceFlag, _ = mockCmd.Flags().GetBoolSlice(boolSliceVarName)
	cmdBoolSliceFlag, _ = cmd.Flags().GetBoolSlice(boolSliceVarName)
	assert.Equal(t, mockBoolSliceFlag, cmdBoolSliceFlag)
	mockBoolSliceFlag, _ = mockCmd.Flags().GetBoolSlice(boolSlicePName)
	cmdBoolSliceFlag, _ = cmd.Flags().GetBoolSlice(boolSlicePName)
	assert.Equal(t, mockBoolSliceFlag, cmdBoolSliceFlag)
	mockBoolSliceFlag, _ = mockCmd.Flags().GetBoolSlice(boolSliceVarPName)
	cmdBoolSliceFlag, _ = cmd.Flags().GetBoolSlice(boolSliceVarPName)
	assert.Equal(t, mockBoolSliceFlag, cmdBoolSliceFlag)
	// string flags
	mockStringFlag, _ := mockCmd.Flags().GetString(stringName)
	cmdStringFlag, _ := cmd.Flags().GetString(stringName)
	assert.Equal(t, mockStringFlag, cmdStringFlag)
	mockStringFlag, _ = mockCmd.Flags().GetString(stringVarName)
	cmdStringFlag, _ = cmd.Flags().GetString(stringVarName)
	assert.Equal(t, mockStringFlag, cmdStringFlag)
	mockStringFlag, _ = mockCmd.Flags().GetString(stringPName)
	cmdStringFlag, _ = cmd.Flags().GetString(stringPName)
	assert.Equal(t, mockStringFlag, cmdStringFlag)
	mockStringFlag, _ = mockCmd.Flags().GetString(stringVarPName)
	cmdStringFlag, _ = cmd.Flags().GetString(stringVarPName)
	assert.Equal(t, mockStringFlag, cmdStringFlag)
	// string slice flags
	mockStringSliceFlag, _ := mockCmd.Flags().GetStringSlice(stringSliceName)
	cmdStringSliceFlag, _ := cmd.Flags().GetStringSlice(stringSliceName)
	assert.Equal(t, mockStringSliceFlag, cmdStringSliceFlag)
	mockStringSliceFlag, _ = mockCmd.Flags().GetStringSlice(stringSliceVarName)
	cmdStringSliceFlag, _ = cmd.Flags().GetStringSlice(stringSliceVarName)
	assert.Equal(t, mockStringSliceFlag, cmdStringSliceFlag)
	mockStringSliceFlag, _ = mockCmd.Flags().GetStringSlice(stringSlicePName)
	cmdStringSliceFlag, _ = cmd.Flags().GetStringSlice(stringSlicePName)
	assert.Equal(t, mockStringSliceFlag, cmdStringSliceFlag)
	mockStringSliceFlag, _ = mockCmd.Flags().GetStringSlice(stringSliceVarPName)
	cmdStringSliceFlag, _ = cmd.Flags().GetStringSlice(stringSliceVarPName)
	assert.Equal(t, mockStringSliceFlag, cmdStringSliceFlag)
}

func getFuncName(function any) string {
	return runtime.FuncForPC(reflect.ValueOf(function).Pointer()).Name()
}
