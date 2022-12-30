/*
 *   sonic-ios-bridge  Connect to your iOS Devices.
 *   Copyright (C) 2022 SonicCloudOrg
 *
 *   This program is free software: you can redistribute it and/or modify
 *   it under the terms of the GNU Affero General Public License as published
 *   by the Free Software Foundation, either version 3 of the License, or
 *   (at your option) any later version.
 *
 *   This program is distributed in the hope that it will be useful,
 *   but WITHOUT ANY WARRANTY; without even the implied warranty of
 *   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *   GNU Affero General Public License for more details.
 *
 *   You should have received a copy of the GNU Affero General Public License
 *   along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */
package run

import (
	"fmt"
	giDevice "github.com/SonicCloudOrg/sonic-gidevice"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/entity"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"net"
	"os"
	"os/signal"
	"regexp"
	"strings"
)

var wdaCmd = &cobra.Command{
	Use:   "wda",
	Short: "Run WebDriverAgent on your devices",
	Long:  `Run WebDriverAgent on your devices`,
	RunE: func(cmd *cobra.Command, args []string) error {
		device := util.GetDeviceByUdId(udid)
		if device == nil {
			os.Exit(0)
		}
		if !strings.HasSuffix(wdaBundleID, ".xctrunner") {
			wdaBundleID += ".xctrunner"
		}
		appList, errList := device.InstallationProxyBrowse(
			giDevice.WithApplicationType(giDevice.ApplicationTypeUser),
			giDevice.WithReturnAttributes("CFBundleShortVersionString", "CFBundleVersion", "CFBundleDisplayName", "CFBundleIdentifier"))
		if errList != nil {
			return util.NewErrorPrint(util.ErrSendCommand, "app list", errList)
		}
		var hasWda = false
		re, _ := regexp.Compile(strings.ReplaceAll(wdaBundleID, "*", "/*"))
		for _, d := range appList {
			a := entity.Application{}
			mapstructure.Decode(d, &a)
			if re.MatchString(a.CFBundleIdentifier) {
				wdaBundleID = a.CFBundleIdentifier
				hasWda = true
				break
			}
		}
		if !hasWda {
			fmt.Printf("%s is not in your device!", wdaBundleID)
			os.Exit(0)
		}
		testEnv := make(map[string]interface{})
		testEnv["USE_PORT"] = serverRemotePort
		testEnv["MJPEG_SERVER_PORT"] = mjpegRemotePort
		output, stopTest, err2 := device.XCTest(wdaBundleID, giDevice.WithXCTestEnv(testEnv))
		if err2 != nil {
			fmt.Printf("WebDriverAgent server start failed: %s", err2)
			os.Exit(0)
		}
		serverListener, err := net.Listen("tcp", fmt.Sprintf(":%d", serverLocalPort))
		if err != nil {
			return err
		}
		defer serverListener.Close()
		go util.StartProxy()(serverListener, serverRemotePort, device)

		if !disableMjpegProxy {
			mjpegListener, err := net.Listen("tcp", fmt.Sprintf(":%d", mjpegLocalPort))
			if err != nil {
				return err
			}
			defer mjpegListener.Close()
			go util.StartProxy()(mjpegListener, mjpegRemotePort, device)
		}

		shutWdaDown := make(chan os.Signal, 1)
		signal.Notify(shutWdaDown, os.Interrupt, os.Kill)

		if !disableShowLog {
			go func() {
				for {
					select {
					case s, ok := <-output:
						if ok {
							fmt.Print(s)
							if strings.Contains(s, "ServerURLHere->") {
								fmt.Println("WebDriverAgent server start successful")
							}
						} else {
							return
						}
					case <-shutWdaDown:
						return
					}
				}
				shutWdaDown <- os.Interrupt
			}()
		}

		<-shutWdaDown
		stopTest()
		fmt.Println("stopped")

		return nil
	},
}

var (
	wdaBundleID       string
	serverRemotePort  int
	mjpegRemotePort   int
	serverLocalPort   int
	mjpegLocalPort    int
	disableMjpegProxy bool
	disableShowLog    bool
)

func initWda() {
	runRootCMD.AddCommand(wdaCmd)
	wdaCmd.Flags().BoolVarP(&disableMjpegProxy, "disable-mjpeg-proxy", "", false, "disable mjpeg-server proxy")
	wdaCmd.Flags().BoolVarP(&disableShowLog, "disable-show-log", "", false, "disable print wda logs")
	wdaCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber ( default first device )")
	wdaCmd.Flags().StringVarP(&wdaBundleID, "bundleId", "b", "com.facebook.WebDriverAgentRunner.xctrunner", "WebDriverAgentRunner bundleId")
	wdaCmd.Flags().IntVarP(&serverRemotePort, "server-remote-port", "", 8100, "WebDriverAgentRunner server remote port")
	wdaCmd.Flags().IntVarP(&mjpegRemotePort, "mjpeg-remote-port", "", 9100, "mjpeg-server remote port")
	wdaCmd.Flags().IntVarP(&serverLocalPort, "server-local-port", "", 8100, "WebDriverAgentRunner server local port")
	wdaCmd.Flags().IntVarP(&mjpegLocalPort, "mjpeg-local-port", "", 9100, "mjpeg-server local port")
}
