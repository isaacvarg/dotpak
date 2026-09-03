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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
