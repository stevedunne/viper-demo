package settings

import (
	"viper-demo/go-viper"

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

	viper.BindFlag(flagset, "log-level", "", "l", "info", "Log level (debug, info, warn, error)")
	viper.BindFlag(flagset, "port", "server", "p", "8080", "Server port")

}
