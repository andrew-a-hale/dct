package version

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "dev"

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display DCT version",
	Long:  `Print the current version number of the DCT tool`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("dct file checker " + version)
	},
}
