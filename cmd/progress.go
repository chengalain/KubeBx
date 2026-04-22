package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/cheng-alain/kubebx/internal/exercises"
	"github.com/cheng-alain/kubebx/internal/k8s"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var progressCmd = &cobra.Command{
	Use:   "progress",
	Short: "Show your learning progress",
	Long:  "Display which exercises you've completed and which are in progress.",
	Run: func(cmd *cobra.Command, args []string) {
		// Load all exercises
		exs, err := exercises.LoadAll()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading exercises: %v\n", err)
			os.Exit(1)
		}

		// Get Kubernetes client to check active exercises
		clientset, err := k8s.GetClient()
		if err != nil {
			fmt.Println("Progress Overview")
			fmt.Println("(Cluster not running - start with 'kbx init')")
			displayExerciseList(exs, nil, nil)
			return
		}

		// List all namespaces to build active/completed maps
		namespaces, err := clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
		activeMap := make(map[string]bool)
		completedMap := make(map[string]bool)
		if err == nil {
			for _, ns := range namespaces.Items {
				if len(ns.Name) >= 5 && ns.Name[:4] == "kbx-" {
					id := ns.Name[4:]
					activeMap[id] = true
					if ns.Labels["kubebx/completed"] == "true" {
						completedMap[id] = true
					}
				}
			}
		}

		fmt.Println("Progress Overview")
		fmt.Println()
		displayExerciseList(exs, activeMap, completedMap)

		// Summary
		total := len(exs)
		completed := len(completedMap)
		active := len(activeMap) - completed

		fmt.Println("\n" + strings.Repeat("─", 50))
		fmt.Printf("Total exercises: %d\n", total)
		fmt.Printf("Completed: %d\n", completed)
		fmt.Printf("In progress: %d\n", active)
		fmt.Printf("Remaining: %d\n", total-len(activeMap))

		if len(activeMap) > 0 {
			fmt.Println("\nNext steps:")
			fmt.Println("   Continue: kbx next")
			fmt.Println("   Check solution: kbx check <id>")
		} else {
			fmt.Println("\nGet started:")
			fmt.Println("   kbx start 01")
		}
	},
}

func displayExerciseList(exs []exercises.Exercise, activeMap map[string]bool, completedMap map[string]bool) {
	for _, ex := range exs {
		status := "[ ]"
		statusText := "Not started"

		if activeMap != nil && activeMap[ex.ID] {
			status = "[>]"
			statusText = "In progress"
		}
		if completedMap != nil && completedMap[ex.ID] {
			status = "[✓]"
			statusText = "Completed"
		}

		typeLabel := fmt.Sprintf("[%s]", ex.Type)
		fmt.Printf("%s %s - %-30s %-10s %s\n", status, ex.ID, ex.Name, typeLabel, statusText)
	}
}

func init() {
	rootCmd.AddCommand(progressCmd)
}
