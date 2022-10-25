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
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
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
			giDevice.WithReturnAttributes("CFBundleVersion", "CFBundleDisplayName", "CFBundleIdentifier"))
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
				kv := strings.Split(value, "=")
				testEnv[kv[0]] = kv[1]
			}
			log.Println("Read env:", testEnv)
		}

		util.CheckMount(device)
		output, stopTest, err2 := device.XCTest(xcTestBundleID, giDevice.WithXCTestEnv(testEnv))
		if err2 != nil {
			fmt.Printf("xctest start failed: %s", err2)
			os.Exit(0)
		}

		shutXcTestDown := make(chan os.Signal, syscall.SIGTERM)
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
