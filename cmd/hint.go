package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cheng-alain/kubebx/internal/exercises"
	"github.com/spf13/cobra"
)

var hintCmd = &cobra.Command{
	Use:   "hint <exercise-id>",
	Short: "Show a hint for an exercise",
	Long:  "Display helpful hints and tips for completing an exercise.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		exerciseID := args[0]

		// Find exercise
		ex, err := exercises.FindByID(exerciseID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Check if hint.md exists
		hintPath := filepath.Join(ex.Path, "hint.md")
		content, err := os.ReadFile(hintPath)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("No hints available for exercise %s yet.\n", ex.ID)
				fmt.Println("\nTry reading the README:")
				fmt.Printf("  cat %s\n", filepath.Join(ex.Path, "README.md"))
				return
			}
			fmt.Fprintf(os.Stderr, "Error reading hint: %v\n", err)
			os.Exit(1)
		}

		// Display hint
		fmt.Println(string(content))
	},
}

func init() {
	rootCmd.AddCommand(hintCmd)
}
