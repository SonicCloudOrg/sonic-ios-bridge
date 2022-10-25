/*
 *  Copyright (C) [SonicCloudOrg] Sonic Project
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 *
 */
package run

import (
	"fmt"
	giDevice "github.com/SonicCloudOrg/sonic-gidevice"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/entity"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"
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
			giDevice.WithReturnAttributes("CFBundleVersion", "CFBundleDisplayName", "CFBundleIdentifier"))
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
		util.CheckMount(device)
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
		go proxy()(serverListener, serverRemotePort, device)

		if !disableMjpegProxy {
			mjpegListener, err := net.Listen("tcp", fmt.Sprintf(":%d", mjpegLocalPort))
			if err != nil {
				return err
			}
			defer mjpegListener.Close()
			go proxy()(mjpegListener, mjpegRemotePort, device)
		}

		shutWdaDown := make(chan os.Signal, syscall.SIGTERM)
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

func proxy() func(listener net.Listener, port int, device giDevice.Device) {
	return func(listener net.Listener, port int, device giDevice.Device) {
		for {
			var accept net.Conn
			var err error
			if accept, err = listener.Accept(); err != nil {
				log.Println("accept:", err)
			}
			fmt.Println("accept", accept.RemoteAddr())
			rInnerConn, err := device.NewConnect(port)
			if err != nil {
				fmt.Println("connect to device fail")
				os.Exit(0)
			}
			rConn := rInnerConn.RawConn()
			rConn.SetDeadline(time.Time{})
			go func(lConn net.Conn) {
				go func(lConn, rConn net.Conn) {
					if _, err := io.Copy(lConn, rConn); err != nil {
						log.Println("local -> remote failed:", err)
					}
				}(lConn, rConn)
				go func(lConn, rConn net.Conn) {
					if _, err := io.Copy(rConn, lConn); err != nil {
						log.Println("local <- remote failed:", err)
					}
				}(lConn, rConn)
			}(accept)
		}
	}
}
