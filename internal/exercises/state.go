package exercises

import (
	"context"
	"fmt"
	"strconv"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// GetCurrentExercise finds the current (last started) exercise by checking namespaces
func GetCurrentExercise(clientset *kubernetes.Clientset) (string, error) {
	currentID, _, err := GetCurrentExerciseWithContext(clientset)
	return currentID, err
}

// GetCurrentExerciseWithContext returns the current exercise and all active exercises
func GetCurrentExerciseWithContext(clientset *kubernetes.Clientset) (string, []string, error) {
	namespaces, err := clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return "", nil, fmt.Errorf("failed to list namespaces: %w", err)
	}

	maxID := 0
	activeExercises := []string{}

	for _, ns := range namespaces.Items {
		// Skip namespaces that are terminating
		if ns.Status.Phase == corev1.NamespaceTerminating {
			continue
		}

		if len(ns.Name) >= 5 && ns.Name[:4] == "kbx-" {
			idStr := ns.Name[4:]
			id, err := strconv.Atoi(idStr)
			if err != nil {
				continue
			}

			activeExercises = append(activeExercises, idStr)

			if id > maxID {
				maxID = id
			}
		}
	}

	if maxID == 0 {
		return "", nil, fmt.Errorf("no active exercise found. Start one with: kbx start 01")
	}

	return fmt.Sprintf("%02d", maxID), activeExercises, nil
}

// GetNextExercise returns the next exercise after the given ID
func GetNextExercise(currentID string) (*Exercise, error) {
	current, err := strconv.Atoi(currentID)
	if err != nil {
		return nil, fmt.Errorf("invalid exercise ID: %s", currentID)
	}

	nextID := fmt.Sprintf("%02d", current+1)

	ex, err := FindByID(nextID)
	if err != nil {
		return nil, fmt.Errorf("no next exercise found. You've completed all available exercises!")
	}

	return ex, nil
}
