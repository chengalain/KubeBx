package cmd

import (
	"fmt"
	"os"

	"github.com/cheng-alain/kubebx/internal/exercises"
	"github.com/cheng-alain/kubebx/internal/k8s"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check <exercise-id>",
	Short: "Check if an exercise is completed correctly",
	Long:  "Verify your solution using the Kubernetes API to validate all requirements.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		exerciseID := args[0]

		// Find exercise
		ex, err := exercises.FindByID(exerciseID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("🔍 Checking exercise %s: %s\n\n", ex.ID, ex.Name)

		// Get Kubernetes client
		clientset, err := k8s.GetClient()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Failed to connect to cluster: %v\n", err)
			fmt.Println("\nMake sure your cluster is running:")
			fmt.Println("  kbx init")
			os.Exit(1)
		}

		// Get checker for this exercise
		checker, err := exercises.NewChecker(exerciseID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ %v\n", err)
			os.Exit(1)
		}

		// Run checks
		result, err := checker.Check(clientset)
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error during check: %v\n", err)
			os.Exit(1)
		}

		// Display results
		if len(result.Details) > 0 {
			for _, detail := range result.Details {
				fmt.Println(detail)
			}
			fmt.Println()
		}

		if len(result.Failures) > 0 {
			for _, failure := range result.Failures {
				fmt.Println(failure)
			}
			fmt.Println()
		}

		// Final message
		if result.Success {
			fmt.Println("✅", result.Message)
			fmt.Println("\nCongratulations! You can move to the next exercise:")
			fmt.Println("  kbx next")
		} else {
			fmt.Println("❌", result.Message)
			fmt.Println("\nNeed help?")
			fmt.Printf("  kbx hint %s\n", ex.ID)
		}
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
