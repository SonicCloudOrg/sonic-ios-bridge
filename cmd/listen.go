package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/conn"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/tool"
	giDevice "github.com/electricbubble/gidevice"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "listener for devices status",
	RunE: func(cmd *cobra.Command, args []string) error {
		if isDetail && (!isJson && !isFormat) {
			return errors.New("detail flag must use with json flag or format flag")
		}
		usbMuxClient, err := giDevice.NewUsbmux()
		if err != nil {
			return tool.NewErrorPrint(tool.ErrConnect, "usbMux", err)
		}
		model := make(chan giDevice.Device)
		shutDownFun, err2 := usbMuxClient.Listen(model)
		if err2 != nil {
			return tool.NewErrorPrint(tool.ErrSendCommand, "listen", err2)
		}

		shutDown := make(chan os.Signal, syscall.SIGTERM)
		signal.Notify(shutDown, os.Interrupt, os.Kill)
		var deviceIdMap = make(map[int]string)
		for {
			select {
			case d := <-model:
				deviceByte, _ := json.Marshal(d.Properties())
				device := &conn.Device{}
				json.Unmarshal(deviceByte, device)
				if len(device.SerialNumber) > 0 {
					deviceIdMap[device.DeviceID] = device.SerialNumber
				} else {
					device.SerialNumber = deviceIdMap[device.DeviceID]
					delete(deviceIdMap, device.DeviceID)
				}
				device.Status = device.GetStatus()
				if device.Status == "online" && isDetail {
					detail, err1 := conn.GetDetail(d)
					if err1 != nil {
						continue
					}
					device.DeviceDetail = *detail
				}
				data := tool.Data(device)
				fmt.Println(tool.Format(data, isFormat, isJson))
			case <-shutDown:
				shutDownFun()
				return nil
			}
		}
		return nil
	},
}

func init() {
	devicesCmd.AddCommand(listenCmd)
	listenCmd.Flags().BoolVarP(&isJson, "json", "j", false, "output for json")
	listenCmd.Flags().BoolVarP(&isFormat, "format", "f", false, "output for json and format")
	listenCmd.Flags().BoolVarP(&isDetail, "detail", "d", false, "output every device's detail, use with json flag or format flag")
}
