package boa

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCobraCmd struct {
	mock.Mock
	*cobra.Command
}

func TestCobraCmdBuilder(t *testing.T) {
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
	mockCmd.AddCommand(subCommands...)
	mockCmd.SetUsageTemplate(usageTemplate)
	mockCmd.SetHelpTemplate(helpTemplate)
	mockCmd.SetUsageFunc(usageFunc)
	mockCmd.SetHelpFunc(helpFunc)

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
}

func getFuncName(function any) string {
	return runtime.FuncForPC(reflect.ValueOf(function).Pointer()).Name()
}
