package k8s

import (
	"fmt"
	"os/exec"

	"github.com/cheng-alain/kubebx/internal/cluster"
)

// ApplyManifest applies a YAML manifest using kubectl
func ApplyManifest(manifestPath string) error {
	kubectlCmd := cluster.GetKubectlCommand()
	cmd := exec.Command(kubectlCmd, "apply", "-f", manifestPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("kubectl apply failed: %w\n%s", err, string(output))
	}

	fmt.Print(string(output))
	return nil
}
