package cmd

import (
	"github.com/spf13/cobra"
)

var bonjourCmd = &cobra.Command{
	Use:   "bonjour",
	Short: "k12-booter bonjour / DNS-ServiceDiscovery tools",
	Long:  ``,
}

func init() {
	rootCmd.AddCommand(bonjourCmd)
}
