package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/isaacvarg/dotpak/internal/manifest"
	"github.com/spf13/cobra"
)

var (
	group       string
	installType string
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all entries in the manifest",
	RunE: func(cmd *cobra.Command, args []string) error {
		hasInstallType := installType != ""
		if hasInstallType {
			it := manifest.InstallType(installType)
			if !it.IsValid() {
				return fmt.Errorf("invalid install type: %q", installType)
			}

		}

		var filtered []manifest.ManifestEntry
		entries, err := manifest.Load()
		if err != nil {
			return err
		}

		for _, e := range entries.Manifest {
			if e.Group != group {
				continue
			}
			if hasInstallType && e.InstallType != manifest.InstallType(installType) {
				continue
			}
			filtered = append(filtered, e)
		}

		w := tabwriter.NewWriter(os.Stdout, 1, 2, 3, ' ', 0)
		fmt.Fprintf(w, "Name\tGroup\tType\t\n")

		for _, e := range filtered {
			fmt.Fprintf(w, "%s\t%s\t%s\n", e.Name, e.Group, e.InstallType)
		}

		err = w.Flush()
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&group, "group", "g", "all", "filter by group")
	listCmd.Flags().StringVarP(&installType, "installType", "i", "", "filter by install type")
}
