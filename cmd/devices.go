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
		usb, err := conn.NewUsbMuxClient()
		if err != nil {
			fmt.Errorf("devices called fail : %w", err)
		}
		list, _ := usb.ListDevices()
		if isJson {
			fmt.Println(list.ToJson())
		} else {
			fmt.Println(list.ToString())
		}
	},
}

func init() {
	rootCmd.AddCommand(devicesCmd)
	devicesCmd.Flags().BoolVarP(&isJson, "json", "j", false, "output format json")
	devicesCmd.Flags().BoolVarP(&isDetail, "detail", "d", false, "output every device's detail")
}
