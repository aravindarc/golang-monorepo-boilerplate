package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang-monorepo-boilerplate/core/log"
	"os"
)

const (
	DefaultConfigFileLocation = "gmb"
	EnvConfigFile             = "GMB_CONFIG"
	FlagConfigFile            = "Config"
	ViperKeyLogLevel          = "log_level"
	ViperKeyPort              = "port"
	ViperKeyDB                = "db"
)

type (
	configDependencies interface {
		log.Provider
	}
	Config struct {
		d configDependencies
		v *viper.Viper
		f *pflag.FlagSet
	}
	Provider interface {
		Config() *Config
	}
)

func NewConfig(d configDependencies, cmd *cobra.Command) *Config {
	v := viper.New()
	f := cmd.Flags()

	slice, _ := f.GetStringSlice("Config")
	for _, s := range slice {
		if _, err := os.Stat(s); err == nil {
			v.SetConfigFile(s)
			goto FINAL
		}
	}
	if envConfigFile, exists := os.LookupEnv(EnvConfigFile); exists && envConfigFile != "" {
		v.SetConfigName(envConfigFile)
	} else {
		v.SetConfigName(DefaultConfigFileLocation)
	}

FINAL:
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	err := v.ReadInConfig()
	if err != nil {
		panic(errors.WithStack(err))
	}

	return &Config{d: d, v: v, f: cmd.Flags()}
}

func (c *Config) LogLevel() string {
	level := c.v.GetString(ViperKeyLogLevel)
	if level == "" {
		return "info"
	}
	return level
}

func (c *Config) Port() string {
	flagPort, _ := c.f.GetString(ViperKeyPort)
	if flagPort != "" {
		return flagPort
	}
	port := c.v.GetString(ViperKeyPort)
	if port == "" {
		return "9876"
	}
	return port
}

func (c *Config) DB() string {
	flagDB, _ := c.f.GetString(ViperKeyDB)
	if flagDB != "" {
		return flagDB
	}
	db := c.v.GetString(ViperKeyDB)
	if db == "" {
		return "gmb.db"
	}
	return db
}

// RegisterServeFlags for serve command
func RegisterServeFlags(flags *pflag.FlagSet) {
	flags.StringSliceP(FlagConfigFile, FlagConfigFile[:1], []string{}, "Config file")
	flags.StringP(ViperKeyPort, ViperKeyPort[:1], "", "port")
	flags.String(ViperKeyDB, "", "db")
}

// RegisterMigrateFlags for migrate command
func RegisterMigrateFlags(flags *pflag.FlagSet) {
	flags.String(ViperKeyDB, "", "db")
}
