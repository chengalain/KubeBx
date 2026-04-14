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
	cleanAll bool
)

var cleanCmd = &cobra.Command{
	Use:   "clean [exercise-id]",
	Short: "Clean up exercise resources",
	Long:  "Remove all Kubernetes resources created for an exercise.",
	Run: func(cmd *cobra.Command, args []string) {
		if cleanAll {
			cleanAllExercises()
			return
		}

		if len(args) == 0 {
			fmt.Fprintln(os.Stderr, "Error: exercise-id required (or use --all flag)")
			fmt.Println("\nUsage:")
			fmt.Println("  kbx clean 01        # Clean exercise 01")
			fmt.Println("  kbx clean --all     # Clean all exercises")
			os.Exit(1)
		}

		exerciseID := args[0]

		// Find exercise
		ex, err := exercises.FindByID(exerciseID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Cleaning exercise %s: %s\n", ex.ID, ex.Name)

		// Get Kubernetes client
		clientset, err := k8s.GetClient()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Failed to connect to cluster: %v\n", err)
			os.Exit(1)
		}

		// Delete namespace
		namespace := fmt.Sprintf("kbx-%s", ex.ID)
		err = clientset.CoreV1().Namespaces().Delete(context.Background(), namespace, metav1.DeleteOptions{})
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Failed to delete namespace: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✓ Deleted namespace %s\n", namespace)
		fmt.Println("\n✅ Exercise cleaned successfully!")
	},
}

func cleanAllExercises() {
	fmt.Println("Cleaning all exercises...")

	// Get Kubernetes client
	clientset, err := k8s.GetClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to connect to cluster: %v\n", err)
		os.Exit(1)
	}

	// List all namespaces starting with kbx-
	namespaces, err := clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to list namespaces: %v\n", err)
		os.Exit(1)
	}

	count := 0
	for _, ns := range namespaces.Items {
		if len(ns.Name) >= 4 && ns.Name[:4] == "kbx-" {
			err := clientset.CoreV1().Namespaces().Delete(context.Background(), ns.Name, metav1.DeleteOptions{})
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to delete namespace %s: %v\n", ns.Name, err)
				continue
			}
			fmt.Printf("✓ Deleted namespace %s\n", ns.Name)
			count++
		}
	}

	if count == 0 {
		fmt.Println("\nNo exercise namespaces found.")
	} else {
		fmt.Printf("\n✅ Cleaned %d exercise(s) successfully!\n", count)
	}
}

func init() {
	rootCmd.AddCommand(cleanCmd)
	cleanCmd.Flags().BoolVar(&cleanAll, "all", false, "Clean all exercises")
}
