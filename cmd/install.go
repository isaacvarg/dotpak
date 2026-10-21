package cmd

import (
	"github.com/isaacvarg/dotpak/internal/install"
	"github.com/spf13/cobra"
)

var installGroup string

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "installs packages from manifest",
	RunE: func(cmd *cobra.Command, args []string) error {
		return install.Install(installGroup)
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().StringVarP(&installGroup, "group", "g", "all", "filter by group")
	//listCmd.Flags().StringVarP(&installType, "installType", "i", "", "filter by install type")
}
