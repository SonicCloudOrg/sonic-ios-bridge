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
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
)

var xctestCmd = &cobra.Command{
	Use:   "xctest",
	Short: "Run xctest on your devices",
	Long:  "Run xctest on your devices",
	RunE: func(cmd *cobra.Command, args []string) error {
		device := util.GetDeviceByUdId(udid)
		if device == nil {
			os.Exit(0)
		}
		appList, errList := device.InstallationProxyBrowse(
			giDevice.WithApplicationType(giDevice.ApplicationTypeUser),
			giDevice.WithReturnAttributes("CFBundleShortVersionString", "CFBundleVersion", "CFBundleDisplayName", "CFBundleIdentifier"))
		if errList != nil {
			return util.NewErrorPrint(util.ErrSendCommand, "app list", errList)
		}
		var hasApp = false
		re, _ := regexp.Compile(strings.ReplaceAll(xcTestBundleID, "*", "/*"))
		for _, d := range appList {
			a := entity.Application{}
			mapstructure.Decode(d, &a)
			if re.MatchString(a.CFBundleIdentifier) {
				xcTestBundleID = a.CFBundleIdentifier
				hasApp = true
				break
			}
		}
		if !hasApp {
			fmt.Printf("%s is not in your device!", xcTestBundleID)
			os.Exit(0)
		}
		testEnv := make(map[string]interface{})
		if len(env) != 0 {
			for _, value := range env {
				if strings.Contains(value, "=") {
					k := value[0:strings.Index(value, "=")]
					v := value[strings.Index(value, "=")+1:]
					testEnv[k] = v
				}
			}
			log.Println("Read env:", testEnv)
		}

		output, stopTest, err2 := device.XCTest(xcTestBundleID, giDevice.WithXCTestEnv(testEnv))
		if err2 != nil {
			fmt.Printf("xctest start failed: %s", err2)
			os.Exit(0)
		}

		shutXcTestDown := make(chan os.Signal, 1)
		signal.Notify(shutXcTestDown, os.Interrupt, os.Kill)

		go func() {
			for s := range output {
				fmt.Print(s)
			}
			shutXcTestDown <- os.Interrupt
		}()

		<-shutXcTestDown
		stopTest()
		fmt.Println("stopped")

		return nil
	},
}

var (
	xcTestBundleID string
	env            []string
)

func initXcTest() {
	runRootCMD.AddCommand(xctestCmd)
	xctestCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber ( default first device )")
	xctestCmd.Flags().StringVarP(&xcTestBundleID, "bundleId", "b", "", "application bundleId")
	xctestCmd.Flags().StringArrayVarP(&env, "env", "e", nil, "test environment params")
}
