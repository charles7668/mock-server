package cmd

import (
	"github.com/spf13/cobra"
	"mock-server/src"
	"mock-server/src/models"
)

func Run() {
	var launchOptions models.LaunchOptions

	var rootCmd = &cobra.Command{
		Use:   "app",
		Short: "A mock backend for frontend testing.",
		Run: func(cmd *cobra.Command, args []string) {
			src.Start(launchOptions)
		},
	}

	rootCmd.Flags().StringVarP(&launchOptions.File, "file", "f", "config.json", "config file path")

	rootCmd.Execute()
}
