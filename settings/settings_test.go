package settings

import (
	"os"
	"path/filepath"
	"testing"
	viperhelper "viper-demo/viper-helper"

	"github.com/spf13/pflag"
)

func TestDefaultSettings(t *testing.T) {

	//Arrange
	expectedLevel := "info"
	expectedPort := 8080

	dir := t.TempDir()
	t.Logf("Temp dir: %s", dir)
	t.Chdir(dir)

	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	BindSettingsFlags(flags)

	if err := flags.Parse(nil); err != nil {
		t.Fatal(err)
	}

	// act
	config, err := viperhelper.GetSettings(&Settings{})

	// assert
	if err != nil {
		t.Fatalf("GetSettings returned error: %v", err)
	}
	assertEqual(t, config.LogLevel, expectedLevel, "LogLevel")
	assertEqual(t, config.Server.Port, expectedPort, "Server.Port")
}

func assertEqual(t *testing.T, got, want interface{}, msg string) {
	if got != want {
		t.Errorf("%s: got %v, want %v", msg, got, want)
	}
}

func TestConfigAndEnvSet_EnvironmentTakesPrecedence(t *testing.T) {
	config := "log-level: config\nserver:\n  port: 9002\n"
	env := map[string]string{"VIPERDEMO_LOG_LEVEL": "environment", "VIPERDEMO_SERVER_PORT": "9001"}
	wantLevel := "environment"
	wantPort := 9001

	dir := t.TempDir()
	t.Chdir(dir)

	// Create config file
	if err := os.WriteFile(filepath.Join(dir, "config.yaml"), []byte(config), 0600); err != nil {
		t.Fatal(err)
	}

	// Add env vars
	for envKey, value := range env {
		t.Setenv(envKey, value)
	}

	// add flags
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	BindSettingsFlags(flags)
	if err := flags.Parse(nil); err != nil {
		t.Fatal(err)
	}

	// act
	// get settings and check values
	settings, err := viperhelper.GetSettings(&Settings{})

	// assert
	if err != nil {
		t.Fatalf("GetSettings returned error: %v", err)
	}
	assertEqual(t, settings.LogLevel, wantLevel, "LogLevel")
	assertEqual(t, settings.Server.Port, wantPort, "Server.Port")

}

func TestConfigFileValues(t *testing.T) {
	config := "log-level: config\nserver:\n  port: 9002\n"
	wantLevel := "config"
	wantPort := 9002

	dir := t.TempDir()
	t.Chdir(dir)

	// Create config file
	if err := os.WriteFile(filepath.Join(dir, "config.yaml"), []byte(config), 0600); err != nil {
		t.Fatal(err)
	}

	// add flags
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	BindSettingsFlags(flags)
	if err := flags.Parse(nil); err != nil {
		t.Fatal(err)
	}

	// act
	// get settings and check values
	settings, err := viperhelper.GetSettings(&Settings{})

	// assert
	if err != nil {
		t.Fatalf("GetSettings returned error: %v", err)
	}
	assertEqual(t, settings.LogLevel, wantLevel, "LogLevel")
	assertEqual(t, settings.Server.Port, wantPort, "Server.Port")
}

func TestAllSettingsConfigured_CommandLineTakesPrecedence(t *testing.T) {
	// arrange
	config := "log-level: config\nserver:\n  port: 9002\n"
	env := map[string]string{"VIPERDEMO_LOG_LEVEL": "environment", "VIPERDEMO_SERVER_PORT": "9001"}
	args := []string{"--log-level=command-line", "--port=9003"}
	wantLevel := "command-line"
	wantPort := 9003

	dir := t.TempDir()
	t.Chdir(dir)

	// Create config file
	if err := os.WriteFile(filepath.Join(dir, "config.yaml"), []byte(config), 0600); err != nil {
		t.Fatal(err)
	}

	// Add env vars
	for envKey, value := range env {
		t.Setenv(envKey, value)
	}

	// add flags
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	BindSettingsFlags(flags)
	if err := flags.Parse(args); err != nil {
		t.Fatal(err)
	}

	// act
	// get settings and check values
	settings, err := viperhelper.GetSettings(&Settings{})

	// assert
	if err != nil {
		t.Fatalf("GetSettings returned error: %v", err)
	}
	assertEqual(t, settings.LogLevel, wantLevel, "LogLevel")
	assertEqual(t, settings.Server.Port, wantPort, "Server.Port")

}
