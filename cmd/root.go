package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kbx",
	Short: "KubeBx - Learn Kubernetes through interactive exercises",
	Long: `KubeBx (kbx) — short for KubeBuilderX — is an open source CLI tool to learn Kubernetes
through hands-on exercises on your local machine.

Build your Kubernetes experience step by step with practical exercises
mixing building resources and debugging broken environments.`,
	Version: "0.1.0",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
