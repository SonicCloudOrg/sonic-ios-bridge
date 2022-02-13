package cmd

import (
	"fmt"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/conn"

	"github.com/spf13/cobra"
)

var devicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "Get iOS device list",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("devices called")
		usb, err := conn.NewUsbMuxClient()
		if err != nil {
			fmt.Errorf("devices called fail : %w", err)
		}
		list, _ := usb.ListDevices()
		fmt.Println(list.ToString())
	},
}

func init() {
	rootCmd.AddCommand(devicesCmd)
	proxyCmd.Flags().BoolP("json", "j", false, "target port/unix path")
}
