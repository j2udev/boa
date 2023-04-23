# Boa

[![GoReportCard](https://goreportcard.com/badge/github.com/j2udev/boa)](https://goreportcard.com/report/github.com/j2udev/boa)
[![Go Reference](https://pkg.go.dev/badge/github.com/j2udev/boa.svg)](https://pkg.go.dev/github.com/j2udev/boa)

Boa is a wrapper for the popular [Cobra](https://github.com/spf13/cobra)
library. It streamlines the building of Cobra Commands, making them easier to
create, read and maintain.

If you initialize a new [cobra-cli](https://github.com/spf13/cobra-cli) project,
you'll end up with something like this:

```go
// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "example",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.test.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
```

While this is small and manageable at first, things can quickly get messy.
Conversely, this is (roughly) the same command using boa:

```go
func NewRootCmd() *cobra.Command {
	long := `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`

	return boa.NewCobraCmd("root").
		WithShortDescription("A brief description of your application").
		WithLongDescription(long).
		WithSubCommands(NewChildCmd()).
		WithBoolPFlag("toggle", "t", false, "Help message for toggle").
		WithStringPPersistentFlag("verbosity", "V", "info", "How verbose should command output be").
		Build()
}
```

where `NewChildCmd()` is defined an another file

```go
func NewChildCmd() *cobra.Command {
	return boa.NewCobraCmd("child").
		WithShortDescription("A subcommand of root").
		WithLongDescription("A detailed description for the child command").
		WithArgs(cobra.MatchAll(cobra.MinimumNArgs(1), cobra.OnlyValidArgs)).
		WithValidArgs([]string{"arg1", "arg2"}).
		WithRunFunc(childFunc).
		Build()
}

func childFunc(cmd *cobra.Command, args []string) {
	//business logic here; recommend abstracting it to a separate package that is cobra agnostic
}
```

and your main package is kept as minimal as possible

```go
func main() {
	err := cmd.NewRootCmd().Execute()
	if err != nil {
		log.Fatal(err)
	}
}
```
