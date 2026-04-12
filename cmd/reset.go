package cmd

import (
	"fmt"
	"os"

	"github.com/cheng-alain/kubebx/internal/cluster"
	"github.com/spf13/cobra"
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Delete the KubeBx cluster completely",
	Long: `Destroy the Kind cluster and remove all KubeBx data.
This will delete:
  - The Kind cluster 'kubebx'
  - All exercise namespaces and resources
  
You will need to run 'kbx init' again to start fresh.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("⚠️  This will destroy the KubeBx cluster and all exercise data.")
		fmt.Print("Continue? (y/N): ")

		var response string
		fmt.Scanln(&response)

		if response != "y" && response != "Y" {
			fmt.Println("Cancelled.")
			return
		}

		fmt.Println("\n🗑️  Deleting cluster...")

		// Check if cluster exists
		exists, err := cluster.ClusterExists()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error checking cluster: %v\n", err)
			os.Exit(1)
		}

		if !exists {
			fmt.Println("No KubeBx cluster found.")
			return
		}

		// Delete cluster
		if err := cluster.DeleteCluster(); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Failed to delete cluster: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("\n✅ Cluster deleted successfully!")
		fmt.Println("\nTo start fresh, run: kbx init")
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)
}
