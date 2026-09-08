package cmd

import (
	"fmt"

	"github.com/isaacvarg/dotpak/internal/groups"
	"github.com/spf13/cobra"
)

var groupCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Creates a package group",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		g, err := groups.Load()
		if err != nil {
			return fmt.Errorf("reading groups: %w", err)
		}

		if err := g.Add(name); err != nil {
			return fmt.Errorf("creating group %q: %w", name, err)
		}

		if err := g.Save(); err != nil {
			return fmt.Errorf("saving groups: %w", err)
		}

		fmt.Printf("Created group %q\n", name)
		return nil
	},
}

func init() {
	groupCmd.AddCommand(groupCreateCmd)
}
