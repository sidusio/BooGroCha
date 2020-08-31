package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func VersionCmd(appName, version string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Book Grouprooms at Chalmers: %s:%s\n", appName, version)
		},
	}
}
