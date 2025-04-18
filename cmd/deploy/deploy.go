package deploy

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	environment string
	force       bool
)

// DeployCmd represents the deploy command
var DeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy your application to different environments",
	Long: `Deploy command allows you to deploy your application to different environments
such as development, staging, or production.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Deploying to %s environment\n", environment)
		if force {
			fmt.Println("Force flag is set, bypassing pre-deployment checks")
		}
		fmt.Println("Deployment started...")
	},
}

func init() {
	DeployCmd.Flags().StringVarP(&environment, "environment", "e", "development", "Environment to deploy to (development, staging, production)")
	DeployCmd.Flags().BoolVarP(&force, "force", "f", false, "Force deployment bypassing pre-deployment checks")
}
