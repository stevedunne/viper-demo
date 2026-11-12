package main

import (
	"fmt"
	"os"
	"viper-demo/settings"
	viperhelper "viper-demo/viper-helper"

	"github.com/spf13/cobra"
)

func main() {

	var command = &cobra.Command{
		Use:   "settings-demo",
		Short: "Show how to use settings from config files, env vars and flags",
		Run: func(cmd *cobra.Command, args []string) {

			settings, err := viperhelper.GetSettings(&settings.Settings{})

			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Settings: %+v\n", settings)
		},
	}

	settings.BindSettingsFlags(command.Flags())

	if err := command.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
