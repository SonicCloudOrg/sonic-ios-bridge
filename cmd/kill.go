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
package cmd

import (
	"fmt"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/entity"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	giDevice "github.com/electricbubble/gidevice"
	"github.com/mitchellh/mapstructure"
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
		result, errList := device.InstallationProxyBrowse(
			giDevice.WithBundleIDs(bundleId),
			giDevice.WithReturnAttributes("CFBundleVersion", "CFBundleDisplayName", "CFBundleIdentifier"))
		if errList != nil {
			return util.NewErrorPrint(util.ErrSendCommand, "appList", errList)
		}
		if result == nil {
			fmt.Printf("%s is not in your device!", bundleId)
			os.Exit(0)
		}
		processList, errProcess := device.AppRunningProcesses()
		if errProcess != nil {
			return util.NewErrorPrint(util.ErrSendCommand, "processList", errList)
		}
		a := entity.Application{}
		mapstructure.Decode(result, &a)
		for _, process := range processList {

		}
		return nil
	},
}

func init() {
	appCmd.AddCommand(killCmd)
	killCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber")
	killCmd.Flags().StringVarP(&bundleId, "bundleId", "b", "", "target bundleId")
	killCmd.MarkFlagRequired("bundleId")
}
