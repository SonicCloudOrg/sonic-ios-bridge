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
