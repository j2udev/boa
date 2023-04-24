package boa

import (
	"log"
	"os"

	"github.com/adrg/xdg"
	"github.com/spf13/viper"
)

// ViperCfgBuilder is a builder that wraps viper Viper objects to allow more
// fluently defining configuration.
type ViperCfgBuilder struct {
	cfg *viper.Viper
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

// Read will read in the config based on configured file/path/name/type.
func (b *ViperCfgBuilder) Read() *ViperCfgBuilder {
	b.cfg.ReadInConfig()
	return b
}

// Build returns a viper Viper from a ViperCfgBuilder
func (b *ViperCfgBuilder) Build() *viper.Viper {
	return b.cfg
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
