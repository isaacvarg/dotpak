package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "prints the dotpak version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("dotpak version %s\n", version)
	},
}

func init() {
	rootCmd.Version = version
	rootCmd.AddCommand(versionCmd)
}
