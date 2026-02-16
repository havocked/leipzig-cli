package cmd

import (
	"fmt"

	"github.com/havocked/leipzig-cli/internal/source/leipzigde"
	"github.com/spf13/cobra"
)

var sourcesCmd = &cobra.Command{
	Use:   "sources",
	Short: "List available event sources",
	Run: func(cmd *cobra.Command, args []string) {
		sources := []struct {
			id, status, desc string
		}{
			{leipzigde.New().ID(), "enabled", "City of Leipzig official event calendar"},
		}
		fmt.Printf("%-15s %-10s %s\n", "SOURCE", "STATUS", "DESCRIPTION")
		for _, s := range sources {
			fmt.Printf("%-15s %-10s %s\n", s.id, s.status, s.desc)
		}
	},
}

func init() {
	rootCmd.AddCommand(sourcesCmd)
}
