# Viper settings demo

This Cobra application demonstrates loading settings from command-line flags, a YAML config file, environment variables, and defaults. 

`main.go` creates and executes the command, 
`settings.go` defines the settings and their flags
`viper.go` handles Viper configuration, config-file loading, environment bindings, and unmarshalling.

## Precedence

Settings use this order, from highest to lowest priority:

1. A supplied command-line flag, such as `--log-level` or `--port`.
2. An environment variable.
3. The corresponding value in the YAML config file.
4. The built-in default.

The sample config is discovered as `./config/config.yaml` or `./config.yaml`. 

`log-level` defaults to `info`, and `server.port` defaults to `8080`.

Environment variables use the `VIPERDEMO_` prefix, with dots and hyphens converted to underscores. For example, `log-level` uses `VIPERDEMO_LOG_LEVEL`, and `server.port` uses `VIPERDEMO_SERVER_PORT`. The environment prefix can be configured in the constants section of [viper.go](viper.go).

## Run

Basic run - will use the config file values, or fall back to the defaults if no config file is present 
```sh
go run .
```

Here the defaults and config values will be overridden by the values specified in the command line 
```sh
go run . --log-level debug --port 9000
```

Use an environment variable to specify the log level
```sh
VIPERDEMO_LOG_LEVEL=warn go run .
```

Here we specify an environment variable and a command line value so the command line will take precidence 
```sh
VIPERDEMO_LOG_LEVEL=warn go run . -l info
```

## Test 

```sh
go test ./...
```

The tests cover defaults, config values, environment-over-config precedence, and command-line overrides.