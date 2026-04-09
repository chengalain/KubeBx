package cmd

import (
	"fmt"
	"os"

	"github.com/cheng-alain/kubebx/internal/exercises"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all available exercises",
	Long:  "Display all available Kubernetes exercises with their type (build/debug) and status.",
	Run: func(cmd *cobra.Command, args []string) {
		exs, err := exercises.LoadAll()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading exercises: %v\n", err)
			os.Exit(1)
		}

		if len(exs) == 0 {
			fmt.Println("No exercises found. Create your first exercise in the exercises/ directory.")
			return
		}

		fmt.Println("Available exercises:\n")

		for _, ex := range exs {
			typeLabel := fmt.Sprintf("[%s]", ex.Type)
			fmt.Printf("%s - %-30s %s\n", ex.ID, ex.Name, typeLabel)
		}

		fmt.Println("\nUse 'kbx start <number>' to begin an exercise")
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
