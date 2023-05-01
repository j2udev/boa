package boa

import (
	"io"
	"log"
	"os"
	"strings"

	"github.com/adrg/xdg"
	"github.com/spf13/viper"
)

// ViperCfgBuilder is a builder that wraps viper.Viper objects to allow more
// fluently defining configuration.
type ViperCfgBuilder struct {
	cfg *viper.Viper
}

// ToViperCfgBuilder is used to convert a viper.Viper object to a
// ViperCfgBuilder
func ToViperCfgBuilder(cmd *viper.Viper) *ViperCfgBuilder {
	return &ViperCfgBuilder{cmd}
}

// NewViperCfg initializes a new viper instance and returns a builder.
func NewViperCfg() *ViperCfgBuilder {
	return &ViperCfgBuilder{
		cfg: viper.New(),
	}
}

// NewDefaultViperCfg initializes a new viper instance with default
// configuration and returns a builder. It adds the user's current working
// directory and XDG_CONFIG_HOME to the searchable config path in that
// respective order and searches for configuration files of 'name' and any
// extension.
func NewDefaultViperCfg(name string) *ViperCfgBuilder {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	b := &ViperCfgBuilder{
		cfg: viper.New(),
	}
	b.cfg.AddConfigPath(cwd)
	b.cfg.AddConfigPath(xdg.ConfigHome + "/" + name)
	b.cfg.SetConfigName(name)
	b.cfg.ReadInConfig()
	return b
}

// WithConfigFiles takes a variable number of filepaths to check for viper
// configuration. The order of the files passed is the order of precedence
// given to each filepath.
func (b *ViperCfgBuilder) WithConfigFiles(files ...string) *ViperCfgBuilder {
	for _, f := range files {
		if exists(f) {
			b.cfg.SetConfigFile(f)
			break
		}
	}
	return b
}

// WithConfigPaths adds a variable number of paths for Viper to search for the
// config file in. It will only add the path if it exists.
func (b *ViperCfgBuilder) WithConfigPaths(paths ...string) *ViperCfgBuilder {
	for _, p := range paths {
		if exists(p) {
			b.cfg.AddConfigPath(p)
		}
	}
	return b
}

// WithConfigName sets the config name to search for in the configured paths.
func (b *ViperCfgBuilder) WithConfigName(name string) *ViperCfgBuilder {
	b.cfg.SetConfigName(name)
	return b
}

// WithConfigType sets the file extension type to search for on config files.
// e.g. "json"
func (b *ViperCfgBuilder) WithConfigType(ext string) *ViperCfgBuilder {
	b.cfg.SetConfigType(ext)
	return b
}

// WithEnvPrefix sets the prefix to use for subsequent bound env vars.
func (b *ViperCfgBuilder) WithEnvPrefix(prefix string) *ViperCfgBuilder {
	b.cfg.SetEnvPrefix(prefix)
	return b
}

// WithBoundEnv binds a Viper key to a ENV variable. ENV variables are case sensitive.
// If only a key is provided, it will use the env key matching the key, uppercased.
// If more arguments are provided, they will represent the env variable names that
// should bind to this key and will be taken in the specified order.
// EnvPrefix will be used when set when env name is not provided.
func (b *ViperCfgBuilder) WithBoundEnv(input ...string) *ViperCfgBuilder {
	b.cfg.BindEnv(input...)
	return b
}

// WithAutomaticEnv makes Viper check if environment variables match any of the existing keys
// (config, default or flags). If matching env vars are found, they are loaded into Viper.
func (b *ViperCfgBuilder) WithAutomaticEnv() *ViperCfgBuilder {
	b.cfg.AutomaticEnv()
	return b
}

// WithEnvKeyReplacer sets the strings.Replacer on the viper object
// Useful for mapping an environmental variable to a key that does
// not match it.
func (b *ViperCfgBuilder) WithEnvKeyReplacer(replacer *strings.Replacer) *ViperCfgBuilder {
	b.cfg.SetEnvKeyReplacer(replacer)
	return b
}

// WithDefaultEnvKeyReplacer is commonly used to convert underscore delimited
// environment variables to their dot delimited equivalents.
//
// For example, the env var COMMAND_NAME can be referenced through Viper as
// viper.Get("command.name")
func (b *ViperCfgBuilder) WithDefaultEnvKeyReplacer() *ViperCfgBuilder {
	b.cfg.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	return b
}

// ReadConfig will read a configuration file, setting existing keys to nil if the
// key does not exist in the file.
//
// If an error is encountered, logs fatal
func (b *ViperCfgBuilder) ReadConfig(in io.Reader) *ViperCfgBuilder {
	err := b.cfg.ReadConfig(in)
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}
	return b
}

// ReadInConfig will discover and load the configuration file from disk
// and key/value stores, searching in one of the defined paths.
//
// If an error is encountered, logs fatal
func (b *ViperCfgBuilder) ReadInConfig() *ViperCfgBuilder {
	err := b.cfg.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading in config: %v", err)
	}
	return b
}

// Build returns a viper.Viper object from a ViperCfgBuilder
func (b *ViperCfgBuilder) Build() *viper.Viper {
	return b.cfg
}

// ReadAndBuild will read in the config based on configured file/path/name/type
// and return a viper.Viper object from a ViperCfgBuilder.
//
// If an error is encountered, logs fatal
func (b *ViperCfgBuilder) ReadInConfigAndBuild() *viper.Viper {
	return b.ReadInConfig().Build()
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
