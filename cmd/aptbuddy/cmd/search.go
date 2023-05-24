package cmd

import (
	"schnoddelbotz/k12-booter/aptbuddy"
	"schnoddelbotz/k12-booter/utility"

	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "INCOMPLETE search on apt bleve index. Use bleve CLI...",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		b, err := aptbuddy.Open(language)
		utility.Fatal(err)
		b.Experiments()
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
