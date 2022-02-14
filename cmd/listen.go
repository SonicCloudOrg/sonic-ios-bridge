package cmd

import (
	"fmt"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/conn"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

var listenCmd = &cobra.Command{
	Use:   "listen",
	Short: "listener for devices status",
	RunE: func(cmd *cobra.Command, args []string) error {
		go func() {
			for {
				usbMuxClient, err1 := conn.NewUsbMuxClient()
				defer usbMuxClient.GetDeviceConn().Close()
				if err1 != nil {
					fmt.Errorf("fail to connect usbMux : %w", err1)
					continue
				}
				receiveFun, err2 := usbMuxClient.Listen()
				if err2 != nil {
					fmt.Errorf("fail to send listen : %w", err2)
				}
				for {
					msg, err := receiveFun()
					if err != nil {
						break
					}
					if isJson {
						fmt.Println(msg.ToJson())
					} else {
						fmt.Println(msg.ToString())
					}
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
	listenCmd.Flags().BoolVarP(&isJson, "json", "j", false, "output format json")
}
