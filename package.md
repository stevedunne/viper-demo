# Making the Viper helper importable

The current code is part of `package main`, so another Go application cannot import it as a library. The helper also depends on demo-specific settings, the process-wide Viper singleton, hard-coded config discovery, and CLI-oriented error handling. The following updates are needed to turn it into a reusable package.

## 1. Create a library package

Move the reusable code into a non-`main` directory, for example `settings/` or `config/`, and change its package declaration accordingly:

```text
viper-demo/
  cmd/settings-demo/main.go
  settings/settings.go
  settings/viper.go
```

The library package must not contain `func main`. Keep Cobra command creation and execution in the demo application. Choose a real module path in `go.mod` before publishing; consumers will import the package using that module path, for example `example.com/project/settings`.

Avoid naming the package `viper`: that name is easily confused with the upstream `github.com/spf13/viper` dependency.

## 2. Separate reusable settings from application settings

Move `Settings`, `ServerSettings`, and their mapstructure tags out of `package main` only if they are intended to be shared by consuming applications. If they describe only this demo, keep those structs in the application and make the helper accept a destination supplied by the caller, or provide a generic loading API.

Export the intended API with Go-style names and documentation, such as `Load`, `Config`, `Option`, and `BindFlags`. Keep implementation details such as config-file initialization private unless callers have a concrete need to invoke them directly.

## 3. Replace Viper's package global with an instance

Use `viper.New()` to create a `*viper.Viper` owned by a loader or configuration object. Pass that instance to setup, flag binding, and unmarshalling instead of calling package-level functions such as `viper.SetDefault`, `viper.BindPFlag`, and `viper.Unmarshal`.

This avoids settings leaking between library users and makes multiple independent configurations and reliable tests possible in one process. A loader should own its Viper instance for its full lifetime.

## 4. Make inputs configurable

Do not force consuming applications to use this demo's config name, search paths, environment prefix, or flag set. Provide options or a configuration struct for the values that need to vary:

- Config file path, or config name and search paths, plus config type.
- Environment prefix, key replacer, and optional explicit environment-variable aliases.
- Defaults and setting keys.
- Optional flag set, if the library is responsible for binding command-line flags.

Prefer accepting a fully qualified config-file path when appropriate; relative paths otherwise resolve from the consuming process's working directory, not the library source directory. Retain the precedence expected by Viper: explicitly changed flags, environment, config file, then defaults.

## 5. Return errors instead of controlling the host process

`BindFlag` currently calls `os.Exit` when binding fails. `GetSettings` panics on unmarshal failure, and config discovery prints directly to stdout/stderr. A library must not terminate or write unrequested output on behalf of its caller.

Return errors from operations that can fail, for example:

```go
func Load[T any](opts ...Option) (T, error)
func BindFlag(v *viper.Viper, flags *pflag.FlagSet, ...) error
```

The exact API can vary, but let the application decide whether to log, retry, or exit. Decide and document whether a missing config file is an error or an allowed condition where environment variables and defaults remain usable.

## 6. Keep Cobra optional

The reusable package may depend on `pflag` if flag binding is part of its public API, but it should not require Cobra. An application can pass `command.Flags()` to the helper, while applications using another CLI framework can load configuration without importing Cobra.

## 7. Add library-focused tests and an import example

Move tests with the package and test through its exported API. Cover independent loader instances, custom config paths and environment prefixes, missing-config behavior, malformed config/unmarshal errors, and precedence across changed flags, environment, config, and defaults. Ensure tests do not depend on global Viper state or leave environment variables / working-directory changes behind.

Add a small example in the README showing a separate application importing the package, creating its Cobra or pflag flags, loading settings, and handling returned errors. Run `go test ./...` from the module and, ideally, compile a small consumer package to prove the public API is genuinely importable.
