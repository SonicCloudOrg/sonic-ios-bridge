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
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	"os"

	"github.com/spf13/cobra"
)

var launchCmd = &cobra.Command{
	Use:   "launch",
	Short: "Launch App",
	Long:  "Launch App",
	RunE: func(cmd *cobra.Command, args []string) error {
		device := util.GetDeviceByUdId(udid)
		if device == nil {
			os.Exit(0)
		}
		_, errLaunch := device.AppLaunch(bundleId)
		if errLaunch != nil {
			return util.NewErrorPrint(util.ErrSendCommand, "launch", errLaunch)
		}
		return nil
	},
}

func initAppLaunch() {
	appRootCMD.AddCommand(launchCmd)
	launchCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber")
	launchCmd.Flags().StringVarP(&bundleId, "bundleId", "b", "", "target bundleId")
	launchCmd.MarkFlagRequired("bundleId")
}
