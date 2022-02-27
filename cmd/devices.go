package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/entity"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	giDevice "github.com/electricbubble/gidevice"
	"github.com/spf13/cobra"
)

var devicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "Get iOS device list",
	Long:  "Get iOS device list",
	RunE: func(cmd *cobra.Command, args []string) error {
		if isDetail && (!isJson && !isFormat) {
			return errors.New("detail flag must use with json flag or format flag")
		}
		usbMuxClient, err := giDevice.NewUsbmux()
		if err != nil {
			return util.NewErrorPrint(util.ErrConnect, "usbMux", err)
		}
		list, err1 := usbMuxClient.Devices()
		if err1 != nil {
			return util.NewErrorPrint(util.ErrSendCommand, "listDevices", err1)
		}
		var deviceList entity.DeviceList
		for _, d := range list {
			deviceByte, _ := json.Marshal(d.Properties())
			device := &entity.Device{}
			if isDetail {
				detail, err2 := entity.GetDetail(d)
				if err2 != nil {
					return err2
				}
				device.DeviceDetail = *detail
			}
			json.Unmarshal(deviceByte, device)
			device.Status = device.GetStatus()
			deviceList.DeviceList = append(deviceList.DeviceList, *device)
		}
		data := util.ResultData(deviceList)
		fmt.Println(util.Format(data, isFormat, isJson))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(devicesCmd)
	devicesCmd.Flags().BoolVarP(&isJson, "json", "j", false, "convert to JSON string")
	devicesCmd.Flags().BoolVarP(&isFormat, "format", "f", false, "convert to JSON string and format")
	devicesCmd.Flags().BoolVarP(&isDetail, "detail", "d", false, "output every device's detail, use with json flag or format flag")
}
