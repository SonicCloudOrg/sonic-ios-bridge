package cmd

import (
	"errors"
	"fmt"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/conn"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/tool"
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
		go func() {
			for {
				usbMuxClient, err1 := conn.NewUsbMuxClient()
				defer usbMuxClient.GetDeviceConn().Close()
				if err1 != nil {
					tool.NewErrorPrint(tool.ErrConnect, "usbMux", err1)
					continue
				}
				receiveFun, err2 := usbMuxClient.Listen()
				if err2 != nil {
					tool.NewErrorPrint(tool.ErrSendCommand, "listen", err2)
					break
				}
				for {
					device, err := receiveFun()
					if err != nil {
						break
					}
					if device.Status == "online" && isDetail {
						detail, err1 := device.GetDetail()
						if err1 != nil {
							fmt.Errorf("get udId %s device detail fail : %w", device.Properties.SerialNumber, err1)
							continue
						}
						device.DeviceDetail = *detail
					}
					data := tool.Data(device)
					fmt.Println(tool.Format(data, isFormat, isJson))
				}
			}
		}()
		//syscall.SIGKILL or syscall.SIGTERM ? SIGTERM can use defer
		signalSetting := make(chan os.Signal, syscall.SIGTERM)
		signal.Notify(signalSetting, os.Interrupt)
		<-signalSetting
		return nil
	},
}

func init() {
	devicesCmd.AddCommand(listenCmd)
	listenCmd.Flags().BoolVarP(&isJson, "json", "j", false, "output for json")
	listenCmd.Flags().BoolVarP(&isFormat, "format", "f", false, "output for json and format")
	listenCmd.Flags().BoolVarP(&isDetail, "detail", "d", false, "output every device's detail, use with json flag or format flag")
}
