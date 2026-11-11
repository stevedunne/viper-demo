package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {

	var command = &cobra.Command{
		Use:   "settings-demo",
		Short: "Show how to use settings from config files, env vars and flags",
		Run: func(cmd *cobra.Command, args []string) {

			settings := GetSettings()

			fmt.Printf("Settings: %+v\n", settings)
		},
	}

	BindSettingsFlags(command.Flags())

	if err := command.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
