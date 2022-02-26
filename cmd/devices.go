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
			tool.NewErrorPrint(tool.ErrConnect, "usbMux", err)
		}
		list, _ := usbMuxClient.Devices()
		var deviceList []conn.Device
		for _, d := range list {
			if isDetail {
				detail, err1 := d.GetValue("", "")

				if err1 != nil {
					return fmt.Errorf("get %s device detail fail : %w", d.Properties().SerialNumber, err1)
				}
				data, _ := json.Marshal(detail)
				d1 := &conn.DeviceDetail{}
				json.Unmarshal(data, d1)

				data2, _ := json.Marshal(d.Properties())
				d2 := &conn.Device{DeviceDetail: *d1}
				json.Unmarshal(data2, d2)
				//dresult, _ := json.Marshal(d2)
				deviceList = append(deviceList,*d2)
				//fmt.Println(string(dresult))
			}
		}
		result:=make(map[string]interface{})
		result["deviceList"] = deviceList
		r, _ := json.Marshal(result)
		fmt.Println(string(r))
		//data := tool.Data(list)
		//fmt.Println(tool.Format(data, isFormat, isJson))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(devicesCmd)
	devicesCmd.Flags().BoolVarP(&isJson, "json", "j", false, "output format json")
	devicesCmd.Flags().BoolVarP(&isFormat, "format", "f", false, "output for json and format")
	devicesCmd.Flags().BoolVarP(&isDetail, "detail", "d", false, "output every device's detail, use with json flag or format flag")
}
