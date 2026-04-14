package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/cheng-alain/kubebx/internal/exercises"
	"github.com/cheng-alain/kubebx/internal/k8s"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	cleanPrevious bool
)

var nextCmd = &cobra.Command{
	Use:   "next [current-exercise-id]",
	Short: "Move to the next exercise",
	Long: `Start the next exercise in the sequence.
	
By default, detects your current exercise and starts the next one.
You can optionally specify which exercise to move from.

Examples:
  kbx next           # Auto-detect current and start next
  kbx next 02        # Start exercise after 02 (i.e., 03)
  kbx next --clean   # Clean current exercise before starting next`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get Kubernetes client
		clientset, err := k8s.GetClient()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Failed to connect to cluster: %v\n", err)
			fmt.Println("\nMake sure your cluster is running:")
			fmt.Println("  kbx init")
			os.Exit(1)
		}

		var currentID string

		// Determine current exercise
		if len(args) > 0 {
			currentID = args[0]
		} else {
			// Auto-detect and check for multiple active exercises
			detected, activeExercises, err := exercises.GetCurrentExerciseWithContext(clientset)
			if err != nil {
				fmt.Fprintf(os.Stderr, "❌ %v\n", err)
				os.Exit(1)
			}

			currentID = detected

			// Warn if multiple exercises are active
			if len(activeExercises) > 1 {
				fmt.Printf("⚠️  Multiple exercises active: %s\n", formatExerciseList(activeExercises))
				fmt.Printf("Using the most recent: %s\n\n", currentID)
			}
		}

		// Verify current exercise exists
		current, err := exercises.FindByID(currentID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Current exercise: %s - %s\n", current.ID, current.Name)

		// Clean previous if requested
		if cleanPrevious {
			fmt.Printf("\n🧹 Cleaning exercise %s...\n", current.ID)
			namespace := fmt.Sprintf("kbx-%s", current.ID)
			err := clientset.CoreV1().Namespaces().Delete(context.Background(), namespace, metav1.DeleteOptions{})
			if err != nil {
				fmt.Fprintf(os.Stderr, "⚠️  Warning: failed to clean namespace: %v\n", err)
			} else {
				fmt.Printf("✓ Cleaned namespace %s\n", namespace)
			}
		}

		// Get next exercise
		next, err := exercises.GetNextExercise(currentID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\n%v\n", err)
			fmt.Println("\n💡 Tip: Specify which exercise to move from:")
			fmt.Printf("  kbx next %s\n", currentID)
			return
		}

		fmt.Printf("\n🚀 Starting next exercise: %s - %s\n", next.ID, next.Name)
		fmt.Printf("Type: %s\n\n", next.Type)

		// Start the next exercise (same logic as kbx start)
		setupPath := fmt.Sprintf("%s/setup.yaml", next.Path)
		if _, err := os.Stat(setupPath); err == nil {
			fmt.Println("Applying initial setup...")
			if err := k8s.ApplyManifest(setupPath); err != nil {
				fmt.Fprintf(os.Stderr, "\n❌ Error applying setup: %v\n", err)
				os.Exit(1)
			}
		}

		fmt.Println("\n✅ Exercise environment ready!")
		fmt.Println("\n💡 Using kubectl:")
		fmt.Println("   Use 'kbx kubectl' for all kubectl commands")
		fmt.Printf("   Example: kbx kubectl get pods -n kbx-%s\n", next.ID)
		fmt.Println("\n📖 Next steps:")
		fmt.Printf("   1. Read instructions: cat %s/README.md\n", next.Path)
		fmt.Printf("   2. Work on the exercise using 'kbx kubectl'\n")
		fmt.Printf("   3. Check your solution: kbx check %s\n", next.ID)
		fmt.Printf("   4. Need help? kbx hint %s\n", next.ID)
	},
}

func formatExerciseList(exercises []string) string {
	result := ""
	for i, ex := range exercises {
		if i > 0 {
			result += ", "
		}
		result += ex
	}
	return result
}

func init() {
	rootCmd.AddCommand(nextCmd)
	nextCmd.Flags().BoolVar(&cleanPrevious, "clean", false, "Clean the current exercise before starting the next one")
}
