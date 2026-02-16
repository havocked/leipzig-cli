package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "leipzig",
	Short: "Discover events and activities in Leipzig",
	Long:  "A CLI tool for discovering events in Leipzig from multiple sources.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
