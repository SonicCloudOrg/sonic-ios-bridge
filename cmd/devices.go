package cmd

import (
	"fmt"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/conn"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/tool"

	"github.com/spf13/cobra"
)

var devicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "Get iOS device list",
	RunE: func(cmd *cobra.Command, args []string) error{
		usb, err := conn.NewUsbMuxClient()
		if err != nil {
			tool.NewErrorPrint(tool.ErrConnect, "usbMux", err)
		}
		list, _ := usb.ListDevices()
		d:=list.DeviceList[0]
		c,err:=conn.GetValueFromDevice(d)
		if err != nil {
			return err
		}
		fmt.Println(c)
		data := tool.Data(list)
		if isJson {
			fmt.Println(data.ToJson())
		} else {
			fmt.Println(data.ToString())
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(devicesCmd)
	devicesCmd.Flags().BoolVarP(&isJson, "json", "j", false, "output format json")
	devicesCmd.Flags().BoolVarP(&isDetail, "detail", "d", false, "output every device's detail")
}
