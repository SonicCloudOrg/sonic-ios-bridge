package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/conn"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/tool"
	giDevice "github.com/electricbubble/gidevice"
	"github.com/spf13/cobra"
)

var devicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "Get iOS device list",
	RunE: func(cmd *cobra.Command, args []string) error {
		if isDetail && (!isJson && !isFormat) {
			return errors.New("detail flag must use with json flag or format flag")
		}
		usbMuxClient, err := giDevice.NewUsbmux()
		if err != nil {
			return tool.NewErrorPrint(tool.ErrConnect, "usbMux", err)
		}
		list, err1 := usbMuxClient.Devices()
		if err1 != nil {
			return tool.NewErrorPrint(tool.ErrSendCommand, "listDevices", err1)
		}
		var deviceList conn.DeviceList
		for _, d := range list {
			deviceByte, _ := json.Marshal(d.Properties())
			device := &conn.Device{}
			if isDetail {
				detail, err2 := conn.GetDetail(d)
				if err2 != nil {
					return err2
				}
				device.DeviceDetail = *detail
			}
			json.Unmarshal(deviceByte, device)
			deviceList.DeviceList = append(deviceList.DeviceList, *device)
		}
		data := tool.Data(deviceList)
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
