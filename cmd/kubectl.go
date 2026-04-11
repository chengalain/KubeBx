package cmd

import (
	"os"
	"os/exec"

	"github.com/cheng-alain/kubebx/internal/cluster"
	"github.com/spf13/cobra"
)

var kubectlCmd = &cobra.Command{
	Use:   "kubectl [kubectl args...]",
	Short: "Run kubectl using KubeBx's managed version",
	Long: `Execute kubectl commands using the kubectl binary managed by KubeBx.
This allows you to use kubectl without installing it globally.

Examples:
  kbx kubectl get pods
  kbx kubectl get pods -n kbx-01
  kbx kubectl run my-pod --image=nginx -n kbx-01
  kbx kubectl describe pod my-pod -n kbx-01`,
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		kubectlPath := cluster.GetKubectlCommand()

		kubectlExec := exec.Command(kubectlPath, args...)
		kubectlExec.Stdout = os.Stdout
		kubectlExec.Stderr = os.Stderr
		kubectlExec.Stdin = os.Stdin

		if err := kubectlExec.Run(); err != nil {
			// kubectl already prints its own error messages
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(kubectlCmd)
}
