package cmd

import (
	"os"

	"github.com/charmbracelet/x/term"
	"github.com/isaacvarg/dotpak/internal/tui"
	"github.com/spf13/cobra"
)

var (
	listGroup       string
	listInstallType string
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all entries in the manifest",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := tui.RunList(listGroup, listInstallType, term.IsTerminal(os.Stdout.Fd()), nil)
		return err
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&listGroup, "group", "g", "", "filter by group (default: all entries)")
	listCmd.Flags().StringVarP(&listInstallType, "installType", "i", "", "filter by install type")
}
