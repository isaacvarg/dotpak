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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
