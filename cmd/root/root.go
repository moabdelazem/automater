package root

import (
	"github.com/moabdelazem/automater/cmd/deploy"
	"github.com/moabdelazem/automater/cmd/monitor"
	"github.com/spf13/cobra"
)

// Version will be set during build time
var Version = "dev"

var rootCmd = &cobra.Command{
	Use:     "automater",
	Version: Version,
	Short:   "Automater - A DevOps automation tool",
	Long: `Automater is a CLI tool designed to automate common DevOps tasks.
It simplifies repetitive operations and helps streamline your workflow.`,
	Run: func(cmd *cobra.Command, args []string) {
		// This is the action that will be executed when the command is called without subcommands
		cmd.Help()
	},
}

// Execute executes the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Initialize flags and configurations here
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Add subcommands
	rootCmd.AddCommand(deploy.DeployCmd)
	rootCmd.AddCommand(monitor.MonitorCmd)
}
