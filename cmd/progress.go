package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/cheng-alain/kubebx/internal/exercises"
	"github.com/cheng-alain/kubebx/internal/k8s"
	"github.com/spf13/cobra"
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
			displayExerciseList(exs, nil)
			return
		}

		// Get active exercises
		_, activeExercises, err := exercises.GetCurrentExerciseWithContext(clientset)
		if err != nil && len(activeExercises) == 0 {
			activeExercises = []string{}
		}

		activeMap := make(map[string]bool)
		for _, id := range activeExercises {
			activeMap[id] = true
		}

		fmt.Println("Progress Overview")
		fmt.Println()
		displayExerciseList(exs, activeMap)

		// Summary
		total := len(exs)
		active := len(activeExercises)

		fmt.Println("\n" + strings.Repeat("─", 50))
		fmt.Printf("Total exercises: %d\n", total)
		fmt.Printf("Active: %d\n", active)
		fmt.Printf("Remaining: %d\n", total-active)

		if active > 0 {
			fmt.Println("\nNext steps:")
			fmt.Println("   Continue: kbx next")
			fmt.Println("   Check solution: kbx check <id>")
		} else {
			fmt.Println("\nGet started:")
			fmt.Println("   kbx start 01")
		}
	},
}

func displayExerciseList(exs []exercises.Exercise, activeMap map[string]bool) {
	for _, ex := range exs {
		status := "[ ]"
		statusText := "Not started"

		if activeMap != nil && activeMap[ex.ID] {
			status = "[>]"
			statusText = "In progress"
		}

		typeLabel := fmt.Sprintf("[%s]", ex.Type)
		fmt.Printf("%s %s - %-30s %-10s %s\n", status, ex.ID, ex.Name, typeLabel, statusText)
	}
}

func init() {
	rootCmd.AddCommand(progressCmd)
}
