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
	case "03":
		return &Exercise03Checker{}, nil
	case "04":
		return &Exercise04Checker{}, nil
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
		result.Message = "All checks passed!"
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
		result.Message = "All checks passed!"
	} else {
		result.Message = "Some checks failed"
	}

	return result, nil
}

// Exercise03Checker checks exercise 03 (Service broken)
type Exercise03Checker struct{}

func (c *Exercise03Checker) Check(clientset *kubernetes.Clientset) (*CheckResult, error) {
	result := &CheckResult{
		Success:  true,
		Details:  []string{},
		Failures: []string{},
	}

	namespace := "kbx-03"
	podName := "nginx-pod"
	svcName := "nginx-service"

	// Check pod exists and is running
	pod, err := clientset.CoreV1().Pods(namespace).Get(context.Background(), podName, metav1.GetOptions{})
	if err != nil {
		result.Success = false
		result.Failures = append(result.Failures, fmt.Sprintf("❌ Pod '%s' not found in namespace '%s'", podName, namespace))
		result.Message = "Pod not found"
		return result, nil
	}
	result.Details = append(result.Details, fmt.Sprintf("✓ Pod '%s' exists", podName))

	if pod.Status.Phase != "Running" {
		result.Success = false
		result.Failures = append(result.Failures, fmt.Sprintf("❌ Pod '%s' is in '%s' state, expected 'Running'", podName, pod.Status.Phase))
	} else {
		result.Details = append(result.Details, fmt.Sprintf("✓ Pod '%s' is Running", podName))
	}

	// Check service exists
	svc, err := clientset.CoreV1().Services(namespace).Get(context.Background(), svcName, metav1.GetOptions{})
	if err != nil {
		result.Success = false
		result.Failures = append(result.Failures, fmt.Sprintf("❌ Service '%s' not found in namespace '%s'", svcName, namespace))
		result.Message = "Some checks failed"
		return result, nil
	}
	result.Details = append(result.Details, fmt.Sprintf("✓ Service '%s' exists", svcName))

	// Check selector matches pod labels
	for key, svcVal := range svc.Spec.Selector {
		podVal, exists := pod.Labels[key]
		if !exists {
			result.Success = false
			result.Failures = append(result.Failures, fmt.Sprintf("❌ Service selector '%s=%s' does not match: pod has no label '%s'", key, svcVal, key))
		} else if podVal != svcVal {
			result.Success = false
			result.Failures = append(result.Failures, fmt.Sprintf("❌ Service selector '%s=%s' does not match pod label '%s=%s'", key, svcVal, key, podVal))
		} else {
			result.Details = append(result.Details, fmt.Sprintf("✓ Service selector %s=%s matches pod label", key, svcVal))
		}
	}

	// Check endpoints have at least one ready address
	endpoints, err := clientset.CoreV1().Endpoints(namespace).Get(context.Background(), svcName, metav1.GetOptions{})
	if err != nil {
		result.Success = false
		result.Failures = append(result.Failures, fmt.Sprintf("❌ Could not retrieve endpoints for Service '%s'", svcName))
	} else {
		readyCount := 0
		for _, subset := range endpoints.Subsets {
			readyCount += len(subset.Addresses)
		}
		if readyCount == 0 {
			result.Success = false
			result.Failures = append(result.Failures, fmt.Sprintf("❌ Service '%s' has no ready endpoints — selector may not match any pod", svcName))
		} else {
			result.Details = append(result.Details, fmt.Sprintf("✓ Service '%s' has %d ready endpoint(s)", svcName, readyCount))
		}
	}

	if result.Success {
		result.Message = "All checks passed!"
	} else {
		result.Message = "Some checks failed"
	}

	return result, nil
}

// Exercise04Checker checks exercise 04 (Deployments)
type Exercise04Checker struct{}

func (c *Exercise04Checker) Check(clientset *kubernetes.Clientset) (*CheckResult, error) {
	result := &CheckResult{
		Success:  true,
		Details:  []string{},
		Failures: []string{},
	}

	namespace := "kbx-04"
	deployName := "nginx-deployment"

	// Check deployment exists
	deploy, err := clientset.AppsV1().Deployments(namespace).Get(context.Background(), deployName, metav1.GetOptions{})
	if err != nil {
		result.Success = false
		result.Failures = append(result.Failures, fmt.Sprintf("❌ Deployment '%s' not found in namespace '%s'", deployName, namespace))
		result.Message = "Deployment not created yet"
		return result, nil
	}
	result.Details = append(result.Details, fmt.Sprintf("✓ Deployment '%s' exists", deployName))

	// Check replicas spec
	specReplicas := int32(1)
	if deploy.Spec.Replicas != nil {
		specReplicas = *deploy.Spec.Replicas
	}
	if specReplicas != 3 {
		result.Success = false
		result.Failures = append(result.Failures, fmt.Sprintf("❌ Deployment spec has %d replica(s), expected 3", specReplicas))
	} else {
		result.Details = append(result.Details, "✓ Deployment spec declares 3 replicas")
	}

	// Check ready replicas
	ready := deploy.Status.ReadyReplicas
	if ready < 3 {
		result.Success = false
		result.Failures = append(result.Failures, fmt.Sprintf("❌ Only %d/3 replicas are Ready", ready))
	} else {
		result.Details = append(result.Details, fmt.Sprintf("✓ %d/3 replicas are Ready", ready))
	}

	// Check image
	if len(deploy.Spec.Template.Spec.Containers) == 0 {
		result.Success = false
		result.Failures = append(result.Failures, "❌ Deployment has no containers defined")
	} else {
		image := deploy.Spec.Template.Spec.Containers[0].Image
		if image != "nginx:latest" && image != "nginx" {
			result.Success = false
			result.Failures = append(result.Failures, fmt.Sprintf("❌ Container image is '%s', expected 'nginx:latest'", image))
		} else {
			result.Details = append(result.Details, fmt.Sprintf("✓ Container using image '%s'", image))
		}
	}

	if result.Success {
		result.Message = "All checks passed!"
	} else {
		result.Message = "Some checks failed"
	}

	return result, nil
}
