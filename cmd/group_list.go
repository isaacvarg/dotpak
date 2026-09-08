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
		groups, err := groups.Load()
		if err != nil {
			return err
		}

		for _, g := range groups.Groups {
			fmt.Println(g.Name)
		}
		return nil
	},
}

func init() {
	groupCmd.AddCommand(groupListCmd)
}
