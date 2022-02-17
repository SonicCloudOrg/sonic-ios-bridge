package cmd

import (
	"errors"
	"fmt"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/conn"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/tool"

	"github.com/spf13/cobra"
)

var devicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "Get iOS device list",
	RunE: func(cmd *cobra.Command, args []string) error {
		if isDetail && (!isJson && !isFormat) {
			return errors.New("detail flag must use with json flag or format flag")
		}
		usb, err := conn.NewUsbMuxClient()
		if err != nil {
			tool.NewErrorPrint(tool.ErrConnect, "usbMux", err)
		}
		list, _ := usb.ListDevices()
		if isDetail {
			for i, d := range list.DeviceList {
				detail, err1 := d.GetDetail()
				if err1 != nil {
					return fmt.Errorf("get %s device detail fail : %w", d.Properties.SerialNumber, err1)
				}
				list.DeviceList[i].DeviceDetail = *detail
			}
		}
		data := tool.Data(list)
		fmt.Println(tool.Format(data, isFormat, isJson))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(devicesCmd)
	devicesCmd.Flags().BoolVarP(&isJson, "json", "j", false, "output format json")
	devicesCmd.Flags().BoolVarP(&isFormat, "format", "f", false, "output for json and format")
	devicesCmd.Flags().BoolVarP(&isDetail, "detail", "d", false, "output every device's detail, use with json flag or format flag")
}
