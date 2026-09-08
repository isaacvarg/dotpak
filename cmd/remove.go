package cmd

import (
	"fmt"

	"github.com/isaacvarg/dotpak/internal/manifest"
	"github.com/spf13/cobra"
)

// groupCmd represents the group command
var removeCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "removes package from manifest",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		m, err := manifest.Load()
		if err != nil {
			return fmt.Errorf("reading manifest: %w", err)
		}

		err = m.Remove(name)
		if err != nil {
			return fmt.Errorf("removing entry %q failed: %w", name, err)
		}

		if err := m.Save(); err != nil {
			return fmt.Errorf("saving manifest entry failed: %w", err)
		}

		fmt.Printf("removed entry %s\n", name)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
