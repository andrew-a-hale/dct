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
	Short: "Swiss army knife for data engineers",
	Long:  `DCT provides utilities to quickly inspect, compare, and manipulate flat data files in CSV, JSON, NDJSON, and Parquet formats`,
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
