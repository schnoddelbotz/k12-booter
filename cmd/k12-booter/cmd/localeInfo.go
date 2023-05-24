/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"schnoddelbotz/k12-booter/internationalization"

	"github.com/spf13/cobra"
)

// localeInfoCmd represents the localeInfo command
var localeInfoCmd = &cobra.Command{
	Use:   "localeInfo",
	Short: "WIP I18N debug",
	Long:  `WIP - Disregard`,
	Run: func(cmd *cobra.Command, args []string) {
		info := internationalization.GetLocaleInfo()
		fmt.Printf("Detected locale: %s\n", info.Locale)
		fmt.Printf("Region (~country): %s, confidence: %s\n", info.Region, info.RegionConfidence)
		fmt.Printf("Detected base language: %s, confidence: %s\n", info.Base, info.BaseConfidence)
		fmt.Printf("Detected script: %s, confidence: %s\n", info.Script, info.ScriptConfidence)
		fmt.Printf("CultureInfo(756): %+v\n", internationalization.CultureInfo(756))
		fmt.Printf(`CultureInfo("Australia"): %+v`+"\n", internationalization.CultureInfo("Australia"))
	},
}

func init() {
	rootCmd.AddCommand(localeInfoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// localeInfoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// localeInfoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
