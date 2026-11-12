package viperhelper

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/pflag"
)

type testSettings struct {
	LogLevel string `mapstructure:"log-level"`
	Server   struct {
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`
}

func clearEnvironment(t *testing.T) {
	t.Helper()
	t.Setenv("VIPERDEMO_LOG_LEVEL", "")
	t.Setenv("VIPERDEMO_SERVER_PORT", "")
}

func assertEqual(t *testing.T, got, want interface{}, msg string) {
	t.Helper()
	if got != want {
		t.Errorf("%s: got %v, want %v", msg, got, want)
	}
}

func TestGetSettingsUsesFlagDefaults(t *testing.T) {
	// Arrange

	clearEnvironment(t)
	t.Chdir(t.TempDir())
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	BindFlag(flags, "log-level", "", "l", "info", "Log level")
	BindFlag(flags, "port", "server", "p", "8080", "Server port")
	if err := flags.Parse(nil); err != nil {
		t.Fatal(err)
	}

	// Act
	settings, err := GetSettings(&testSettings{})

	// Assert
	if err != nil {
		t.Fatalf("GetSettings returned error: %v", err)
	}
	assertEqual(t, settings.LogLevel, "info", "LogLevel")
	assertEqual(t, settings.Server.Port, 8080, "Server.Port")
}

func TestGetSettingsReadsConfigFile(t *testing.T) {
	// Arrange

	clearEnvironment(t)
	dir := t.TempDir()
	t.Chdir(dir)
	config := "log-level: config\nserver:\n  port: 9002\n"
	if err := os.WriteFile(filepath.Join(dir, "config.yaml"), []byte(config), 0600); err != nil {
		t.Fatal(err)
	}
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	BindFlag(flags, "log-level", "", "l", "info", "Log level")
	BindFlag(flags, "port", "server", "p", "8080", "Server port")
	if err := flags.Parse(nil); err != nil {
		t.Fatal(err)
	}

	// Act
	settings, err := GetSettings(&testSettings{})

	// Assert
	if err != nil {
		t.Fatalf("GetSettings returned error: %v", err)
	}
	assertEqual(t, settings.LogLevel, "config", "LogLevel")
	assertEqual(t, settings.Server.Port, 9002, "Server.Port")
}

func TestGetSettingsEnvironmentOverridesConfig(t *testing.T) {
	// Arrange
	dir := t.TempDir()
	t.Chdir(dir)
	config := "log-level: config\nserver:\n  port: 9002\n"
	if err := os.WriteFile(filepath.Join(dir, "config.yaml"), []byte(config), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("VIPERDEMO_LOG_LEVEL", "environment")
	t.Setenv("VIPERDEMO_SERVER_PORT", "9001")
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	BindFlag(flags, "log-level", "", "l", "info", "Log level")
	BindFlag(flags, "port", "server", "p", "8080", "Server port")
	if err := flags.Parse(nil); err != nil {
		t.Fatal(err)
	}

	// Act
	settings, err := GetSettings(&testSettings{})

	// Assert
	if err != nil {
		t.Fatalf("GetSettings returned error: %v", err)
	}
	assertEqual(t, settings.LogLevel, "environment", "LogLevel")
	assertEqual(t, settings.Server.Port, 9001, "Server.Port")
}

func TestGetSettingsChangedFlagsOverrideEnvironment(t *testing.T) {
	// Arrange
	dir := t.TempDir()
	t.Chdir(dir)
	config := "log-level: config\nserver:\n  port: 9002\n"
	if err := os.WriteFile(filepath.Join(dir, "config.yaml"), []byte(config), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("VIPERDEMO_LOG_LEVEL", "environment")
	t.Setenv("VIPERDEMO_SERVER_PORT", "9001")
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	BindFlag(flags, "log-level", "", "l", "info", "Log level")
	BindFlag(flags, "port", "server", "p", "8080", "Server port")
	args := []string{"--log-level=command-line", "--port=9003"}
	if err := flags.Parse(args); err != nil {
		t.Fatal(err)
	}

	// Act
	settings, err := GetSettings(&testSettings{})

	// Assert
	if err != nil {
		t.Fatalf("GetSettings returned error: %v", err)
	}
	assertEqual(t, settings.LogLevel, "command-line", "LogLevel")
	assertEqual(t, settings.Server.Port, 9003, "Server.Port")
}

// func TestInitializationCanOnlyhappenOnce(t *testing.T) {

// 	// Arrange

// 	// Act
// 	err1 := InitViperHelper("", "", "", "", "")

// 	err2 := InitViperHelper("", "", "", "", "")

// 	// Assert
// 	if err1 != nil {
// 		t.Fatalf("First InitViperHelper returned error: %v", err1)
// 	}
// 	if err2 == nil {
// 		t.Fatalf("Second InitViperHelper should have returned error but didn't")
// 	}
// }
