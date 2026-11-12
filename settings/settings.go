package settings

import (
	viperhelper "github.com/stevedunne/viper-demo/viper-helper"

	"github.com/spf13/pflag"
)

type Settings struct {
	LogLevel string         `mapstructure:"log-level"`
	Server   ServerSettings `mapstructure:"server"`
}

type ServerSettings struct {
	Port int `json:"port" mapstructure:"port"`
}

func BindSettingsFlags(flagset *pflag.FlagSet) {

	viperhelper.BindFlag(flagset, "log-level", "", "l", "info", "Log level (debug, info, warn, error)")
	viperhelper.BindFlag(flagset, "port", "server", "p", "8080", "Server port")

}
