package main

import "github.com/spf13/pflag"

type Settings struct {
	LogLevel string         `mapstructure:"log-level"`
	Server   ServerSettings `mapstructure:"server"`
}

type ServerSettings struct {
	Port int `json:"port" mapstructure:"port"`
}

func BindSettingsFlags(flagset *pflag.FlagSet) {

	BindFlag(flagset, "log-level", "", "l", "info", "Log level (debug, info, warn, error)")
	BindFlag(flagset, "port", "server", "p", "8080", "Server port")

}
