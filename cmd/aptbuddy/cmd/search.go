package cmd

import (
	"fmt"
	"schnoddelbotz/k12-booter/aptbuddy"
	"schnoddelbotz/k12-booter/utility"

	"github.com/spf13/cobra"
)

var numResults int

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search [flags] QUERY",
	Short: "INCOMPLETE search on apt bleve index. Use bleve CLI...",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		b, err := aptbuddy.Open(language)
		utility.Fatal(err)
		r := b.Search(args[0], numResults, aptbuddy.HighlightDescAndTags)
		fmt.Printf("%+v\n", r)
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)
	searchCmd.Flags().IntVar(&numResults, "numResults", 10, "number of hits to print")
}
