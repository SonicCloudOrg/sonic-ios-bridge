package cmd

import (
	"github.com/SonicCloudOrg/sonic-ios-bridge/cmd/usb"
	"github.com/spf13/cobra"
)

var usbCmd = &cobra.Command{
	Use:   "usb",
	Short: "Control USB hub",
	Long:  "Control USB hub",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.AddCommand(usbCmd)
	usb.InitUsb(usbCmd)
}
