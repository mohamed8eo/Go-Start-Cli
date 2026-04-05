package cmd

import (
	"log"

	tea "charm.land/bubbletea/v2"
	"github.com/mohamed8eo/gostart/internal/tui/add"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Search and install Go packages interactively",
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(add.InitialModel())
		if _, err := p.Run(); err != nil {
			log.Fatalf("error: %s\n", err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
