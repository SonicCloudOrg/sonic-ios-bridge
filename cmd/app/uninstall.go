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
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	"github.com/spf13/cobra"
	"os"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall App in your device",
	Long:  "Uninstall App in your device",
	RunE: func(cmd *cobra.Command, args []string) error {
		device := util.GetDeviceByUdId(udid)
		if device == nil {
			os.Exit(0)
		}
		errUninstall := device.AppUninstall(bundleId)
		if errUninstall != nil {
			return util.NewErrorPrint(util.ErrSendCommand, "uninstall", errUninstall)
		}
		fmt.Println("uninstall successful")
		return nil
	},
}

func initAppUninstall() {
	appRootCMD.AddCommand(uninstallCmd)
	uninstallCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber")
	uninstallCmd.Flags().StringVarP(&bundleId, "bundleId", "b", "", "target bundleId")
	uninstallCmd.MarkFlagRequired("bundleId")
}
