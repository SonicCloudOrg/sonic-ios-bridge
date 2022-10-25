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
package app

import (
	"fmt"
	giDevice "github.com/SonicCloudOrg/sonic-gidevice"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	"os"

	"github.com/spf13/cobra"
)

var killCmd = &cobra.Command{
	Use:   "kill",
	Short: "Kill app process",
	Long:  "Kill app process",
	RunE: func(cmd *cobra.Command, args []string) error {
		device := util.GetDeviceByUdId(udid)
		if device == nil {
			os.Exit(0)
		}
		lookup, err := device.InstallationProxyLookup(giDevice.WithBundleIDs(bundleId))
		if err != nil {
			return util.NewErrorPrint(util.ErrSendCommand, "look up", err)
		}
		lookupResult := lookup.(map[string]interface{})
		if lookupResult[bundleId] == nil {
			fmt.Printf("%s is not in your device!", bundleId)
			os.Exit(0)
		}
		lookupResult = lookupResult[bundleId].(map[string]interface{})
		execName := lookupResult["CFBundleExecutable"]

		processList, errProcess := device.AppRunningProcesses()
		if errProcess != nil {
			return util.NewErrorPrint(util.ErrSendCommand, "processList", errProcess)
		}

		var hit bool
		for _, process := range processList {
			if process.Name == execName {
				hit = true
				device.AppKill(process.Pid)
				fmt.Println("kill process successful!")
				break
			}
		}
		if !hit {
			fmt.Printf("%s process not found\n", bundleId)
		}
		return nil
	},
}

func initAppKill() {
	appRootCMD.AddCommand(killCmd)
	killCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber")
	killCmd.Flags().StringVarP(&bundleId, "bundleId", "b", "", "target bundleId")
	killCmd.MarkFlagRequired("bundleId")
}
