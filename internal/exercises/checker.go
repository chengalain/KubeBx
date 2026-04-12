package exercises

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// CheckResult represents the result of an exercise check
type CheckResult struct {
	Success  bool
	Message  string
	Details  []string
	Failures []string
}

// Checker is the interface for exercise validation
type Checker interface {
	Check(clientset *kubernetes.Clientset) (*CheckResult, error)
}

// NewChecker creates a checker for a specific exercise
func NewChecker(exerciseID string) (Checker, error) {
	switch exerciseID {
	case "01":
		return &Exercise01Checker{}, nil
	case "02":
		return &Exercise02Checker{}, nil
	default:
		return nil, fmt.Errorf("no checker available for exercise %s", exerciseID)
	}
}

// Exercise01Checker checks exercise 01 (Pod basics)
type Exercise01Checker struct{}

func (c *Exercise01Checker) Check(clientset *kubernetes.Clientset) (*CheckResult, error) {
	result := &CheckResult{
		Success:  true,
		Details:  []string{},
		Failures: []string{},
	}

	namespace := "kbx-01"
	podName := "my-first-pod"

	// Check if pod exists
	pod, err := clientset.CoreV1().Pods(namespace).Get(context.Background(), podName, metav1.GetOptions{})
	if err != nil {
		result.Success = false
		result.Failures = append(result.Failures, fmt.Sprintf("❌ Pod '%s' not found in namespace '%s'", podName, namespace))
		result.Message = "Pod not created yet"
		return result, nil
	}

	result.Details = append(result.Details, fmt.Sprintf("✓ Pod '%s' exists in namespace '%s'", podName, namespace))

	// Check if pod is running
	if pod.Status.Phase != "Running" {
		result.Success = false
		result.Failures = append(result.Failures, fmt.Sprintf("❌ Pod is in '%s' state, expected 'Running'", pod.Status.Phase))
	} else {
		result.Details = append(result.Details, "✓ Pod is Running")
	}

	// Check image
	if len(pod.Spec.Containers) == 0 {
		result.Success = false
		result.Failures = append(result.Failures, "❌ Pod has no containers")
	} else {
		image := pod.Spec.Containers[0].Image
		if image != "nginx:latest" && image != "nginx" {
			result.Success = false
			result.Failures = append(result.Failures, fmt.Sprintf("❌ Container image is '%s', expected 'nginx:latest'", image))
		} else {
			result.Details = append(result.Details, fmt.Sprintf("✓ Container using image '%s'", image))
		}
	}

	// Check container ready
	if len(pod.Status.ContainerStatuses) > 0 {
		if pod.Status.ContainerStatuses[0].Ready {
			result.Details = append(result.Details, "✓ Container is ready")
		} else {
			result.Success = false
			result.Failures = append(result.Failures, "❌ Container is not ready")
		}
	}

	if result.Success {
		result.Message = "All checks passed! 🎉"
	} else {
		result.Message = "Some checks failed"
	}

	return result, nil
}

// Exercise02Checker checks exercise 02 (Labels & Selectors)
type Exercise02Checker struct{}

func (c *Exercise02Checker) Check(clientset *kubernetes.Clientset) (*CheckResult, error) {
	result := &CheckResult{
		Success:  true,
		Details:  []string{},
		Failures: []string{},
	}

	namespace := "kbx-02"

	// Expected pods with their labels
	expectedPods := map[string]map[string]string{
		"frontend": {
			"app":  "web",
			"tier": "frontend",
			"env":  "prod",
		},
		"backend": {
			"app":  "api",
			"tier": "backend",
			"env":  "prod",
		},
		"worker": {
			"app":  "processor",
			"tier": "backend",
			"env":  "dev",
		},
	}

	// Check each pod
	for podName, expectedLabels := range expectedPods {
		pod, err := clientset.CoreV1().Pods(namespace).Get(context.Background(), podName, metav1.GetOptions{})
		if err != nil {
			result.Success = false
			result.Failures = append(result.Failures, fmt.Sprintf("❌ Pod '%s' not found in namespace '%s'", podName, namespace))
			continue
		}

		result.Details = append(result.Details, fmt.Sprintf("✓ Pod '%s' exists", podName))

		// Check if pod is running
		if pod.Status.Phase != "Running" {
			result.Success = false
			result.Failures = append(result.Failures, fmt.Sprintf("❌ Pod '%s' is in '%s' state, expected 'Running'", podName, pod.Status.Phase))
		} else {
			result.Details = append(result.Details, fmt.Sprintf("✓ Pod '%s' is Running", podName))
		}

		// Check labels
		for key, expectedValue := range expectedLabels {
			actualValue, exists := pod.Labels[key]
			if !exists {
				result.Success = false
				result.Failures = append(result.Failures, fmt.Sprintf("❌ Pod '%s' missing label '%s'", podName, key))
			} else if actualValue != expectedValue {
				result.Success = false
				result.Failures = append(result.Failures, fmt.Sprintf("❌ Pod '%s' label '%s' is '%s', expected '%s'", podName, key, actualValue, expectedValue))
			} else {
				result.Details = append(result.Details, fmt.Sprintf("✓ Pod '%s' has correct label %s=%s", podName, key, expectedValue))
			}
		}
	}

	if result.Success {
		result.Message = "All checks passed! 🎉"
	} else {
		result.Message = "Some checks failed"
	}

	return result, nil
}
