package cmd

import (
	"schnoddelbotz/k12-booter/aptbuddy"
	"schnoddelbotz/k12-booter/utility"

	"github.com/spf13/cobra"
)

// indexCmd represents the index command
var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Fetch debian package index and index it using bleve",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := aptbuddy.FetchAndIndex(language)
		utility.Fatal(err)
	},
}

func init() {
	rootCmd.AddCommand(indexCmd)
}
