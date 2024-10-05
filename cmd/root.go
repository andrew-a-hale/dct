package cmd

import (
	art "dct/cmd/art"
	"dct/cmd/diff"
	"dct/cmd/peek"
	"dct/cmd/version"
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dct",
	Short: "dct is a data file checking tool",
	Long:  `A data file checking tool to quickly check and compare data files of various formats`,
}

func init() {
	rootCmd.AddCommand(version.VersionCmd)
	rootCmd.AddCommand(art.ArtCmd)
	rootCmd.AddCommand(peek.PeekCmd)
	rootCmd.AddCommand(diff.DiffCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
