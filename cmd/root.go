package cmd

import (
	"dct/cmd/art"
	"dct/cmd/chart"
	"dct/cmd/diff"
	"dct/cmd/flattify"
	"dct/cmd/generator"
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
	rootCmd.AddCommand(chart.ChartCmd)
	rootCmd.AddCommand(peek.PeekCmd)
	rootCmd.AddCommand(diff.DiffCmd)
	rootCmd.AddCommand(generator.GenCmd)
	rootCmd.AddCommand(flattify.FlattifyCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
