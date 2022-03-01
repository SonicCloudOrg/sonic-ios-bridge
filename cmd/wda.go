package cmd

import (
	"fmt"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	giDevice "github.com/electricbubble/gidevice"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var wdaCmd = &cobra.Command{
	Use:   "wda",
	Short: "Run WebDriverAgent on your devices",
	Long:  `Run WebDriverAgent on your devices`,
	RunE: func(cmd *cobra.Command, args []string) error {
		usbMuxClient, err := giDevice.NewUsbmux()
		if err != nil {
			return util.NewErrorPrint(util.ErrConnect, "usbMux", err)
		}
		list, err1 := usbMuxClient.Devices()
		if err1 != nil {
			return util.NewErrorPrint(util.ErrSendCommand, "listDevices", err1)
		}
		if len(list) == 0 {
			fmt.Errorf("no device connected")
			os.Exit(0)
		} else {
			var device giDevice.Device
			if len(udid) != 0 {
				for i, d := range list {
					if d.Properties().SerialNumber == udid {
						device = list[i]
						break
					}
				}
			} else {
				device = list[0]
			}
			if device.Properties().SerialNumber != "" {
				testEnv := make(map[string]interface{})
				testEnv["USE_PORT"] = serverRemotePort
				testEnv["MJPEG_SERVER_PORT"] = mjpegRemotePort
				output, stopTest, err2 := device.XCTest(wdaBundleID, giDevice.WithXCTestEnv(testEnv))
				if err2 != nil {
					fmt.Printf("WebDriverAgent server start failed... try to mount developer disk image...")
					value, err3 := device.GetValue("", "ProductVersion")
					if err3 != nil {
						return util.NewErrorPrint(util.ErrSendCommand, "get value", err3)
					}
					ver := strings.Split(value.(string), ".")
					var reVer string
					if len(ver) >= 2 {
						reVer = ver[0] + "." + ver[1]
					}
					done := util.LoadDevelopImage(reVer)
					if done {
						var dmg = "DeveloperDiskImage.dmg"
						var sign = dmg + ".signature"
						err4 := device.MountDeveloperDiskImage(fmt.Sprintf(".sib/%s/%s", reVer, dmg), fmt.Sprintf(".sib/%s/%s", reVer, sign))
						if err4 != nil {
							fmt.Errorf("mount develop disk image fail")
							os.Exit(0)
						} else {
							output, stopTest, err2 = device.XCTest(wdaBundleID, giDevice.WithXCTestEnv(testEnv))
							if err2 != nil {
								fmt.Errorf("WebDriverAgent server still start failed")
								os.Exit(0)
							}
						}
					} else {
						fmt.Errorf("download develop disk image fail")
						os.Exit(0)
					}
				}
				shutDown := make(chan os.Signal, syscall.SIGTERM)
				signal.Notify(shutDown, os.Interrupt, os.Kill)

				go func() {
					for s := range output {
						fmt.Print(s)
						if strings.Contains(s, "ServerURLHere->") {
							fmt.Println("WebDriverAgent server start successful")
						}
					}
					shutDown <- os.Interrupt
				}()

				<-shutDown
				stopTest()
				fmt.Println("stopped")
			} else {
				fmt.Errorf("device no found")
				os.Exit(0)
			}
		}
		return nil
	},
}

var (
	wdaBundleID      string
	serverRemotePort int
	mjpegRemotePort  int
	serverLocalPort  int
	mjpegLocalPort   int
)

func init() {
	runCmd.AddCommand(wdaCmd)
	wdaCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber ( default first device )")
	wdaCmd.Flags().StringVarP(&wdaBundleID, "bundleId", "b", "com.facebook.WebDriverAgentRunner.xctrunner", "WebDriverAgentRunner bundleId")
	wdaCmd.Flags().IntVarP(&serverRemotePort, "server-remote-port", "", 8100, "WebDriverAgentRunner server remote port")
	wdaCmd.Flags().IntVarP(&mjpegRemotePort, "mjpeg-remote-port", "", 9100, "mjpeg-server remote port")
	wdaCmd.Flags().IntVarP(&serverLocalPort, "server-local-port", "", 8100, "WebDriverAgentRunner server local port")
	wdaCmd.Flags().IntVarP(&mjpegLocalPort, "mjpeg-local-port", "", 9100, "mjpeg-server local port")
}
