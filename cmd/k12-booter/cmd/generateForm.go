/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"schnoddelbotz/k12-booter/formgenerator"
	"schnoddelbotz/k12-booter/utility"

	"github.com/spf13/cobra"
)

// generateFormCmd represents the generateForm command
var generateFormCmd = &cobra.Command{
	Use:   "generateForm",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		inputFilename, err := cmd.Flags().GetString("input-file")
		utility.Fatal(err)
		formgenerator.CreateFormAsNeeded(false, inputFilename)
	},
}

func init() {
	rootCmd.AddCommand(generateFormCmd)

	// Here you will define your flags and configuration settings.
	generateFormCmd.Flags().String("input-file", "", "HTML input file for form generator")
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateFormCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateFormCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
