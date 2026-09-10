package cmd

import (
	"fmt"

	"github.com/isaacvarg/dotpak/internal/groups"
	"github.com/spf13/cobra"
)

var groupListCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all package groups",
	RunE: func(cmd *cobra.Command, args []string) error {
		g, err := groups.Load()
		if err != nil {
			return err
		}

		for _, name := range g.Names() {
			fmt.Println(name)
		}
		return nil
	},
}

func init() {
	groupCmd.AddCommand(groupListCmd)
}
