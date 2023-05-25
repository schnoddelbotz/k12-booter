package cmd

import (
	"schnoddelbotz/k12-booter/dnssd"

	"github.com/spf13/cobra"
)

var browseCmd = &cobra.Command{
	Use:   "browse",
	Short: "browse LAN for k12-booter teacher nodes",
	// it's here for initial development and debugging.
	// the GUI should guide new users to "their" teacher node(s).
	Args: cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		println("Browsing for k12-booter nodes ... Press ctrl-c to quit.")
		dnssd.Browse()
	},
}

func init() {
	bonjourCmd.AddCommand(browseCmd)
}
