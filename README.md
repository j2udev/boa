# Boa

[![Build Status](https://github.com/j2udev/boa/actions/workflows/go.yml/badge.svg)](https://github.com/j2udev/boa/actions/workflows/go.yml)
[![GoReportCard](https://goreportcard.com/badge/github.com/j2udev/boa)](https://goreportcard.com/report/github.com/j2udev/boa)
[![Go Reference](https://pkg.go.dev/badge/github.com/j2udev/boa.svg)](https://pkg.go.dev/github.com/j2udev/boa)

Boa is a wrapper for the popular [Cobra](https://github.com/spf13/cobra) and
[Viper](https://github.com/spf13/viper) libraries. It streamlines the building
of Cobra Commands and Viper configuration, making them easier to create, read
and maintain.

## Disclaimer

This project should be considered unstable until it is officially released. Use
at your own risk.

## CobraCmdBuilder

Boa wraps the construction of Cobra commands in a builder as opposed to the
struct literal approach taken by the
[cobra-cli](https://github.com/spf13/cobra-cli)(and most other Cobra users). If
you initialize a new cobra-cli project, you'll end up with something like this:

<!-- markdownlint-disable MD010 -->

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

<!-- markdownlint-enable MD010 -->

While this is small and manageable at first, things can quickly get messy.
Conversely, this is (roughly) the same command using boa:

<!-- markdownlint-disable MD010 -->

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

<!-- markdownlint-enable MD010 -->

where `NewChildCmd()` is defined in another file

<!-- markdownlint-disable MD010 -->

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

<!-- markdownlint-enable MD010 -->

and your main package is kept as minimal as possible

<!-- markdownlint-disable MD010 -->

```go
func main() {
	err := cmd.NewRootCmd().Execute()
	if err != nil {
		log.Fatal(err)
	}
}
```

<!-- markdownlint-enable MD010 -->

## BoaCmdBuilder

If you are perfectly content with the traditional Cobra CLI in which the
positional args are unknown and therefore aren't listed in the help/usage text,
you likely have no need for a boa.Command. A good example of a CLI like this is
[kubectl](https://kubernetes.io/docs/reference/kubectl/).

```text
kubectl logs -f <what pod name?>
kubectl apply -f <what manifest?>
kubectl describe deployment <what deployment?>
```

In a CLI like kubectl, the CLI doesn't have static information about the
arguments applied to its commands, therefore the default Cobra command works
perfectly fine.

Let's imagine for a second that you're building a CLI that _does_ have static
positional args. For example:

```text
mycoolcli install kubectl helm skaffold
|         |       |       |    |
|         |       |-------|----| static positional args
|         |
|         |- sub command
|
|- root command
```

In a situation like this you likely _do_ want to see the positional args in your
help/usage text. You might also want to see a helpful description about each
argument. To accomplish this, you would ordinarily have to override the
cobra.Command's help/usage function(s)/template(s); however, a boa.Command can
support this without sacrificing the power of the cobra.Command and the
boa.CobraCmdBuilder.

A boa.Command embeds the cobra.Command and wraps it with new fields in an effort
to cover additional uses cases like the one detailed above. Similarly, the
BoaCmdBuilder embeds the CobraCmdBuilder and wraps it with additional methods to
facilitate adding non-cobra.Command native fields and overriding the help/usage
function(s)/template(s) more easily.

A BoaCmdBuilder can seamlessly chain into a CobraCmdBuilder, but not viceversa.
For example, this is valid:

<!-- markdownlint-disable MD010 -->

```go
func NewInstallCmd() *cobra.Command {
	return boa.NewCmd("install").
		WithValidOptions(
			boa.Option{Args: []string{"kubectl"}, Desc: "install kubectl"},
			boa.Option{Args: []string{"helm"}, Desc: "install helm"},
			boa.Option{Args: []string{"skaffold"}, Desc: "install skaffold"},
		).
		WithOptionsTemplate().
		WithMinValidArgs(1).
		WithShortDescription("install tools").
		WithLongDescription("install tools that make a productive kubernetes developer").
		WithRunFunc(install).
		Build()
}
```

<!-- markdownlint-enable MD010 -->

but this is not:

<!-- markdownlint-disable MD010 -->

```go
func NewInstallCmd() *boa.Command {
	return boa.NewCmd("install").
		WithShortDescription("install tools").
		WithLongDescription("install tools that make a productive kubernetes developer").
		WithRunFunc(install).
		WithValidOptions(
			boa.Option{Args: []string{"kubectl"}, Desc: "install kubectl"},
			boa.Option{Args: []string{"helm"}, Desc: "install helm"},
			boa.Option{Args: []string{"skaffold"}, Desc: "install skaffold"},
		).
		WithOptionsTemplate().
		WithMinValidArgs(1).
		Build()
}
```

<!-- markdownlint-enable MD010 -->

`WithValidOptions()`, `WithOptionsTemplate()`, and `WithMinValidArgs()` are
methods on the BoaCmdBuilder which embeds the CobraCmdBuilder and therefore has
access to all of its methods. This is why the BoaCmdBuilder methods can chain
into CobraCmdBuilder methods, but the reverse is not true unless you use the
`ToBoaCmdBuilder()` method on the CobraCmdBuilder. If we look at the previously
invalid example, we can make it valid by chaining `ToBoaCmdBuilder()`.

<!-- markdownlint-disable MD010 -->

```go
func NewInstallCmd() *boa.Command {
	return boa.NewCmd("install").
		WithShortDescription("install tools").
		WithLongDescription("install tools that make a productive kubernetes developer").
		WithRunFunc(install).
		ToBoaCmdBuilder().
		WithValidOptions(
			boa.Option{Args: []string{"kubectl"}, Desc: "install kubectl"},
			boa.Option{Args: []string{"helm"}, Desc: "install helm"},
			boa.Option{Args: []string{"skaffold"}, Desc: "install skaffold"},
		).
		WithOptionsTemplate().
		WithMinValidArgs(1).
		Build()
}
```

<!-- markdownlint-enable MD010 -->

Let's see the help text of `mycoolcli install` now that boa.Command vs
cobra.Command and BoaCmdBuilder vs CobraCmdBuilder have been addressed.

```text
Usage:
  mycoolcli install [flags] [options]

Options:
  kubectl    install kubectl
  helm       install helm
  skaffold   install skaffold

Flags:
  -h, --help   help for install
```

boa.Commands also support profiles, which simply means that you can define an
argument as a shorthand for a set of other options. In practice this looks
something like:

<!-- markdownlint-disable MD010 -->

```go
func NewInstallCmd() *cobra.Command {
	return boa.NewCmd("install").
		WithValidOptions(
			boa.Option{Args: []string{"kubectl"}, Desc: "install kubectl"},
			boa.Option{Args: []string{"helm"}, Desc: "install helm"},
			boa.Option{Args: []string{"skaffold"}, Desc: "install skaffold"},
		).
		WithValidProfiles(
			boa.Profile{Args: []string{"core"}, Opts: []string{"kubectl", "helm"}, Desc: "install core tools for working with k8s"},
			boa.Profile{Args: []string{"developer", "dev"}, Opts: []string{"kubectl", "helm", "skaffold"}, Desc: "install k8s developer tools"},
		).
		WithOptionsTemplate().
		WithMinValidArgs(1).
		WithShortDescription("install tools").
		WithLongDescription("install tools that make a productive kubernetes developer").
		WithRunFunc(install).
		Build()
```

<!-- markdownlint-enable MD010 -->

which results in the following usage:

```text
Usage:
  mycoolcli install [flags] [options]

Options:
  kubectl    install kubectl
  helm       install helm
  skaffold   install skaffold

Profiles:
  core             install core tools for working with k8s
    ↳ Options:     kubectl, helm
  developer, dev   install k8s developer tools
    ↳ Options:     kubectl, helm, skaffold

Flags:
  -h, --help   help for install
```

It's important to note that Options and Profiles simply end up in the
cobra.Commands/boa.Commands `args` slice. It is up to the user to determine what
an option or profile type of argument means for their application.

## ViperCfgBuilder

Currently, Boa doesn't extensively wrap Viper. Viper is moving towards v2 and it
doesn't lend itself to being wrapped in a builder as well as Cobra. That said,
Boa does offer a simple builder for initializing Viper configuration and
includes a sane default configuration that can be used. If your use case is more
complicated than simply pointing at a config file(s) and reading it in, Boa's
`ViperCfgBuilder` probably isn't worth your time. If your use-case _is_ simple,
read on.

To initialize Viper configuration that searches in the user's current working
directory and their XDG_CONFIG_HOME in that respective order, you can use the
`NewDefaultViperCfg()` function.

<!-- markdownlint-disable MD010 -->

```go
// define your configuration schema; viper uses [mapstructure](https://pkg.go.dev/github.com/mitchellh/mapstructure)
// type Schema struct {
// 	Cfg     SomeStruct            `mapstructure:"config"`
// 	MoreCfg map[string]SomeStruct `mapstructure:"moreConfig"`
// }
// var cfg Schema
viper := boa.NewDefaultViperCfg("boa").Build()
err := viper.UnmarshalExact(&cfg)
```

<!-- markdownlint-enable MD010 -->

If the defaults don't work for you, you can always build your own!

<!-- markdownlint-disable MD010 -->

```go
viper := boa.NewViperCfg().
	WithConfigPaths("/potential/path/to/config", "/another/one").
	WithConfigName("my-cool-config").
	ReadInConfig().
	Build()
```

<!-- markdownlint-enable MD010 -->

or

<!-- markdownlint-disable MD010 -->

```go
viper := boa.NewViperCfg().
	WithConfigFiles("/potential/path/to/config.yml", "/another/one/config.json").
	ReadInConfigAndBuild()
```

<!-- markdownlint-enable MD010 -->

It's important to note that the `ReadInConfig()` and `ReadInConfigAndBuild()`
methods can encounter an error and will log fatal if so. If you need to handle
the error differently, `Build()` first, then call `viper.ReadInConfig()`
yourself.
