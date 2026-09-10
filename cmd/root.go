// Package cmd sets up the cli commands for dotpak
package cmd

import (
	"os"

	"github.com/isaacvarg/dotpak/internal/tui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dotpak",
	Short: "Package installation manager for your dotfiles",
	Run: func(cmd *cobra.Command, args []string) {
		tui.App()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
