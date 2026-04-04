/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	tea "charm.land/bubbletea/v2"
	"github.com/mohamed8eo/gostart/internal/tui/tui"
	"github.com/mohamed8eo/gostart/internal/work"
	"github.com/spf13/cobra"
)

var (
	ProjectName string
	FrameWork   string
	DB          string
	SQL         string
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new Go project with optional framework and database setup",
	Long: `Initialize a new Go project scaffold interactively or via CLI flags.
You can choose your project name, preferred Go framework, database, and SQL driver.

Examples:
  # Interactive TUI to choose everything
  gostart init

  # Using CLI flags to set framework and database
  gostart init --framework Gin --database PostgreSQL --sql GORM
`,
	Run: func(cmd *cobra.Command, args []string) {
		// Launch interactive TUI if no CLI args provided
		if len(args) == 0 {
			p := tea.NewProgram(tui.InitialModel())
			result, err := p.Run()
			if err != nil {
				log.Fatalf("error: %s\n", err.Error())
			}

			m, ok := result.(tui.Model)
			if !ok {
				log.Fatalf("Error casting final model")
			}

			// Capture final choices
			ProjectName = m.ProjectName.Value()
			FrameWork = m.Frameworks[m.SelectedFramework].Name
			DB = m.Databases[m.SelectedDatabase].Name
			SQL = m.SqlDrivers[m.SelectedSQL].Name
		} else {
			ProjectName = args[0]
		}
		if err := work.CreateProjectStructure(ProjectName, FrameWork, DB, SQL); err != nil {
			log.Fatalf("Failed to create project: %v", err)
		}
		if err := work.InstallDependencies(ProjectName, FrameWork, DB, SQL); err != nil {
			log.Fatalf("Failed to install dependencies: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// CLI flags
	initCmd.Flags().StringVarP(&FrameWork, "framework", "f", "None", "Select Go framework (Gin, Echo, Fiber, Chi, None)")
	initCmd.Flags().StringVarP(&DB, "database", "d", "None", "Select database driver (PostgreSQL, MySQL, SQLite, MongoDB, None)")
	initCmd.Flags().StringVarP(&SQL, "sql", "s", "None", "Select SQL driver (GORM, sqlx, sqlc, pgx, None)")
}
