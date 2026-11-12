# Viper Settings Demo

This repository contains a Cobra demo application and an importable helper package at `viper-helper/`. The demo's settings type and flag declarations are in `settings/`; the helper binds flags and loads values into a caller-provided struct.

## Import the Package

The module path is `github.com/stevedunne/viper-demo`. Add the package to another application with:

```sh
go get github.com/stevedunne/viper-demo/viper-helper@latest
```

Example using pflag:

```go
package main

import (
	"fmt"
	"os"

	viperhelper "github.com/stevedunne/viper-demo/viper-helper"
)

type Settings struct {
	LogLevel string `mapstructure:"log-level"`
	Server   struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`
}

func main() {

	var command = &cobra.Command{
		Use:   "settings-demo",
		Short: "Show how to use settings from config files, env vars and flags",
		Run: func(cmd *cobra.Command, args []string) {

			settings, err := viperhelper.GetSettings(&settings.Settings{})

			... 
		},
	}

	settings.BindSettingsFlags(command.Flags())

	if err := command.Execute(); err != nil {
		...
	}
}
```

For Cobra, pass `command.Flags()` to `BindFlag` and let Cobra parse the command-line arguments.

## Initialization and Defaults

The helper uses a process-wide singleton. Initialization is optional: without an explicit call, first use selects these defaults:

- Config type: `yaml`
- Config name: `config`
- Config search paths: `./` and `./config/`
- Environment prefix: `VIPERDEMO`
- Key replacer: dots and hyphens become underscores

To configure the singleton, call `InitViperHelper` before binding flags or loading settings. Its arguments are config type, config name, comma-separated search paths, comma-separated key-replacer pairs, and environment prefix. For example:

```go
err := viperhelper.InitViperHelper("yaml", "app", "./,./config/", ".,_,-,_", "MYAPP")
```

Initialization can only happen once. A later call returns an error. `Dispose` clears the singleton; use it only when intentionally resetting the helper, such as in tests.

## Precedence and Errors

Values are resolved in this order: explicitly changed command-line flags, environment variables, config file values, then defaults supplied to `BindFlag`. With the default prefix, `log-level` reads `VIPERDEMO_LOG_LEVEL` and `server.port` reads `VIPERDEMO_SERVER_PORT`.

A missing config file is allowed, so environment variables and defaults can still be used. Other config read failures and unmarshalling failures are returned by `GetSettings`; `BindFlag` also returns binding errors. The application is responsible for handling those errors.

## Run and Test the Demo

```sh
go run .
go run . --log-level debug --port 9000
VIPERDEMO_LOG_LEVEL=warn go run .
VIPERDEMO_LOG_LEVEL=warn go run . -l info
go test ./...
```