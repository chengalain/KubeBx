package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cheng-alain/kubebx/internal/exercises"
	"github.com/cheng-alain/kubebx/internal/k8s"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start <exercise-id>",
	Short: "Start an exercise",
	Long:  "Deploy the initial environment for a Kubernetes exercise.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		exerciseID := args[0]

		// Find exercise
		ex, err := exercises.FindByID(exerciseID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("🚀 Starting exercise %s: %s\n", ex.ID, ex.Name)
		fmt.Printf("Type: %s\n\n", ex.Type)

		// Check if setup.yaml exists
		setupPath := filepath.Join(ex.Path, "setup.yaml")
		if _, err := os.Stat(setupPath); os.IsNotExist(err) {
			fmt.Println("✓ No initial setup required for this exercise")
			fmt.Println("\n📖 Read the instructions:")
			readmePath := filepath.Join(ex.Path, "README.md")
			fmt.Printf("   cat %s\n", readmePath)
			return
		}

		// Apply setup
		fmt.Println("Applying initial setup...")
		if err := k8s.ApplyManifest(setupPath); err != nil {
			fmt.Fprintf(os.Stderr, "\n❌ Error applying setup: %v\n", err)
			fmt.Println("\nMake sure your cluster is running:")
			fmt.Println("  kbx init")
			os.Exit(1)
		}

		fmt.Println("\n✅ Exercise environment ready!")
		fmt.Println("\n💡 Using kubectl:")
		fmt.Println("   Use 'kbx kubectl' for all kubectl commands")
		fmt.Println("   Example: kbx kubectl get pods -n kbx-01")
		fmt.Println("\n📖 Next steps:")
		readmePath := filepath.Join(ex.Path, "README.md")
		fmt.Printf("   1. Read instructions: cat %s\n", readmePath)
		fmt.Printf("   2. Work on the exercise using 'kbx kubectl'\n")
		fmt.Printf("   3. Check your solution: kbx check %s\n", ex.ID)
		fmt.Printf("   4. Need help? kbx hint %s\n", ex.ID)
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
