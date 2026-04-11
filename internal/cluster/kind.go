package cluster

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	ClusterName = "kubebx"
)

// ClusterExists checks if the kubebx cluster already exists
func ClusterExists() (bool, error) {
	kindCmd := GetKindCommand()
	cmd := exec.Command(kindCmd, "get", "clusters")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("failed to list kind clusters: %w", err)
	}

	clusters := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, cluster := range clusters {
		if cluster == ClusterName {
			return true, nil
		}
	}

	return false, nil
}

// CreateCluster creates a new Kind cluster with custom configuration
func CreateCluster() error {
	kindCmd := GetKindCommand()

	// Kind config for the cluster
	config := `kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
name: kubebx
nodes:
- role: control-plane
  extraPortMappings:
  - containerPort: 80
    hostPort: 80
    protocol: TCP
  - containerPort: 443
    hostPort: 443
    protocol: TCP
`

	// Create temp config file
	tmpFile, err := os.CreateTemp("", "kind-config-*.yaml")
	if err != nil {
		return fmt.Errorf("failed to create temp config file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(config); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}
	tmpFile.Close()

	// Create cluster
	fmt.Println("Creating Kind cluster (this may take a few minutes)...")
	cmd := exec.Command(kindCmd, "create", "cluster", "--config", tmpFile.Name())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create cluster: %w", err)
	}

	return nil
}

// DeleteCluster deletes the kubebx Kind cluster
func DeleteCluster() error {
	kindCmd := GetKindCommand()
	cmd := exec.Command(kindCmd, "delete", "cluster", "--name", ClusterName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to delete cluster: %w", err)
	}

	return nil
}

// SetKubeconfig sets the kubectl context to the kubebx cluster
func SetKubeconfig() error {
	kubectlCmd := GetKubectlCommand()
	cmd := exec.Command(kubectlCmd, "config", "use-context", "kind-"+ClusterName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to set context: %w\n%s", err, string(output))
	}

	return nil
}
