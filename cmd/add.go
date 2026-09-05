package cmd

import (
	"fmt"

	"github.com/isaacvarg/dotpak/internal/manifest"
	"github.com/spf13/cobra"
)

var (
	group          string
	installType    string
	installCommand string
)

// groupCmd represents the group command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a package",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		hasInstallCommand := installCommand != ""
		ic := name
		if hasInstallCommand {
			ic = installCommand
		}

		it := manifest.InstallType(installType)
		if !it.IsValid() {
			return fmt.Errorf("invalid install type: %q", installType)
		}

		m, err := manifest.Load()
		if err != nil {
			return fmt.Errorf("reading groups: %w", err)
		}

		err = m.Add(
			name,
			it,
			group,
			ic,
		)
		if err != nil {
			return fmt.Errorf("creating entry %q: %w", name, err)
		}

		if err := m.Save(); err != nil {
			return fmt.Errorf("saving manifest entry: %w", err)
		}


		fmt.Printf("Created entry %q in group %q for %q \n", name, group, installType )

		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(&group, "group", "g", "all", "specifies the group to add this entry to")
	addCmd.Flags().StringVarP(&installType, "installType", "i", "", "type of entry")
	_ = addCmd.MarkFlagRequired("installType")
	addCmd.Flags().StringVarP(&installCommand, "installCommand", "c", "", "simply the name of packages for packages or the full command for something like flatpak")
}
