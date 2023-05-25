package cmd

import (
	"fmt"
	"schnoddelbotz/k12-booter/dnssd"

	"github.com/spf13/cobra"
)

var lookupCmd = &cobra.Command{
	Use:   "lookup",
	Short: "look up k12-booter teacher host via dns-sd/bonjour",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Looking up %s via bonjour...", args[0])
		dnssd.Lookup(args[0]) // ... time out ...?
	},
}

func init() {
	bonjourCmd.AddCommand(lookupCmd)
}
