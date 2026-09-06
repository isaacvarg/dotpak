package cmd

import (
	"fmt"

	"github.com/isaacvarg/dotpak/internal/manifest"
	"github.com/spf13/cobra"
)

var (
	installType string
	group string
	installCommand string
)

// groupCmd represents the group command
var addCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "Add a package, app or plugin to the dotpak manifest",
	Long: `Add a package, app or plugin to the dotpak manifest.

Available install types include:
	pacman
	aur
	flatpak
	omarchy
	mise

The -i, --installType flag allows a install command to be 
specified that is different than the name. For instance, this 
is necessary for flathub appID's or specifying the repo for an 
omarchy plugin. If ommitted, the install command will equal the
name of the package.

`,
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
