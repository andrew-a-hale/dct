package cmd

import (
	"dct/cmd/art"
	"dct/cmd/chart"
	"dct/cmd/diff"
	"dct/cmd/flattify"
	"dct/cmd/generator"
	"dct/cmd/infer"
	"dct/cmd/js2sql"
	"dct/cmd/peek"
	"dct/cmd/profile"
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
	rootCmd.AddCommand(infer.InferCmd)
	rootCmd.AddCommand(diff.DiffCmd)
	rootCmd.AddCommand(generator.GenCmd)
	rootCmd.AddCommand(flattify.FlattifyCmd)
	rootCmd.AddCommand(profile.ProfileCmd)
	rootCmd.AddCommand(js2sql.Js2SqlCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
