package cmd

import (
	"fmt"

	"github.com/isaacvarg/dotpak/internal/manifest"
	"github.com/spf13/cobra"
)

var (
	installType    string
	group          string
	installCommand string
	isPacman       bool
	isAUR          bool
	isFlatpak      bool
	isOmarchy      bool
	isMise         bool
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

Install Command
---------------
The -i or --installCommand flag allows a install command to be specified that is different than the name. For example, the name and command differ for flathub appID's or specifying the repo for an omarchy plugin. If ommitted, the install command will equal the name of the package.

`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// name of the entry
		name := args[0]

		switch {
		case isPacman:
			installType = "pacman"
		case isAUR:
			installType = "aur"
		case isFlatpak:
			installType = "flatpak"
		case isOmarchy:
			installType = "omarchy"
		case isMise:
			installType = "mise"
		}

		if installType == "" {
			return fmt.Errorf("an install type must be specified via -i or --installType or a shortcut flag (-p, -a, -f, -o, -m)")
		}

		// if install command not included this will equal to the name
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
			return fmt.Errorf("reading manifest: %w", err)
		}

		err = m.Add(
			name,
			it,
			group,
			ic,
		)
		if err != nil {
			return fmt.Errorf("creating entry %q failed: %w", name, err)
		}

		if err := m.Save(); err != nil {
			return fmt.Errorf("saving manifest entry: %w", err)
		}

		fmt.Printf("Created entry %q in group %q for %q \n", name, group, installType)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
	addCmd.Flags().StringVarP(&group, "group", "g", "all", "specifies the group to add this entry to")
	addCmd.Flags().StringVarP(&installType, "installType", "i", "", "type of entry")
	addCmd.Flags().StringVarP(&installCommand, "installCommand", "c", "", "simply the name of packages for packages or the full command for something like flatpak")

	// install type boolean flags for shortcuts
	addCmd.Flags().BoolVarP(&isPacman, "pacman", "p", false, "shortcut for pacman install type")
	addCmd.Flags().BoolVarP(&isAUR, "aur", "a", false, "shortcut for aur install type")
	addCmd.Flags().BoolVarP(&isFlatpak, "flatpak", "f", false, "shortcut for flatpak install type")
	addCmd.Flags().BoolVarP(&isOmarchy, "omarchy", "o", false, "shortcut for omarchy install type")
	addCmd.Flags().BoolVarP(&isMise, "mise", "m", false, "shortcut for mise install type")

	// ensure that only one shortcut or the install command provided
	addCmd.MarkFlagsMutuallyExclusive("installType", "pacman", "aur", "flatpak", "omarchy", "mise")
}
