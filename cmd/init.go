package cmd

import (
	"fmt"
	"os"

	"github.com/cheng-alain/kubebx/internal/cluster"
	"github.com/spf13/cobra"
)

var (
	forceRecreate bool
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize the local Kubernetes cluster",
	Long: `Setup a local Kind cluster for running KubeBx exercises.

This command will:
  - Check that Docker is installed and running (required)
  - Download and install Kind if not present
  - Download and install kubectl if not present
  - Create a Kind cluster named 'kubebx'
  - Configure kubectl to use the cluster

All tools except Docker are automatically managed by KubeBx.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Checking and installing dependencies...\n")

		// Check and install dependencies
		if err := cluster.CheckAndInstallDependencies(); err != nil {
			fmt.Fprintf(os.Stderr, "\n❌ %v\n", err)
			os.Exit(1)
		}

		fmt.Println()

		// Check if cluster exists
		exists, err := cluster.ClusterExists()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error checking cluster: %v\n", err)
			os.Exit(1)
		}

		if exists && !forceRecreate {
			fmt.Println("✓ KubeBx cluster already exists")
			fmt.Println("\nCluster is ready! Use 'kbx list' to see available exercises.")
			fmt.Println("\nTo recreate the cluster, use: kbx init --force")
			return
		}

		if exists && forceRecreate {
			fmt.Println("Deleting existing cluster...")
			if err := cluster.DeleteCluster(); err != nil {
				fmt.Fprintf(os.Stderr, "Error deleting cluster: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("✓ Cluster deleted\n")
		}

		// Create cluster
		if err := cluster.CreateCluster(); err != nil {
			fmt.Fprintf(os.Stderr, "\n❌ Failed to create cluster: %v\n", err)
			os.Exit(1)
		}

		// Set kubeconfig context
		if err := cluster.SetKubeconfig(); err != nil {
			fmt.Fprintf(os.Stderr, "\nWarning: %v\n", err)
			fmt.Println("You may need to run: kubectl config use-context kind-kubebx")
		}

		fmt.Println("\n✅ KubeBx cluster is ready!")
		fmt.Println("\nNext steps:")
		fmt.Println("  1. List available exercises: kbx list")
		fmt.Println("  2. Start your first exercise: kbx start 01")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVar(&forceRecreate, "force", false, "Delete and recreate the cluster if it exists")
}
