package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	configFileType      = "yaml"
	configFileLocations = "./,./config/"
	envPrefix           = "VIPERDEMO"
)

func BindFlag(flagset *pflag.FlagSet, name, group, shorthand, def, usage string) {
	//add the flag to the flagset
	flagset.StringP(name, shorthand, def, usage)

	flagName := name
	if group != "" {
		flagName = fmt.Sprintf("%s.%s", group, name)
	}

	// bind the flag to viper
	if err := viper.BindPFlag(flagName, flagset.Lookup(name)); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	// also set the default value in viper, so that if the flag is not set, viper will use the default value
	viper.SetDefault(flagName, def)
}

func GetSettings() *Settings {
	var cfg Settings

	InitViperConfig("config")

	for _, key := range viper.AllKeys() {
		fmt.Printf("Key: %s, Value: %v\n", key, viper.Get(key))
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		// reduce to error once testing is done
		panic(fmt.Errorf("failed to unmarshal config: %w", err))
	}
	return &cfg
}

func InitViperConfig(configFileName string) {
	keyReplacer := strings.NewReplacer(".", "_", "-", "_")
	viper.SetEnvPrefix(envPrefix)
	viper.SetEnvKeyReplacer(keyReplacer)
	viper.AutomaticEnv()

	for _, loc := range strings.Split(configFileLocations, ",") {
		viper.AddConfigPath(loc)
	}
	viper.SetConfigName(configFileName)
	viper.SetConfigType(configFileType)

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Error loading viper: %+v\n", err)
	} else {
		fmt.Fprintf(os.Stderr, "Using config file: %s\n", viper.ConfigFileUsed())
	}

}
