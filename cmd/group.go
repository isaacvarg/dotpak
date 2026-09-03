package cmd

import (
	"github.com/spf13/cobra"
)

// groupCmd represents the group command
var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "Manage package groups",
	Long: `Package groups allow for containers that packages can be added to so that during deployment containers can be conditionally be selected for installation.

An example usage is adding Work and Personal package groups, using dotfile manager such as chezmoi to provide the environment variable and only install Work packages on work deployments.`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(groupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// groupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// groupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
