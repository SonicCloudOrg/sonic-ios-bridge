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
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	"os"

	"github.com/spf13/cobra"
)

var rebootCmd = &cobra.Command{
	Use:   "reboot",
	Short: "Reboot or Shutdown device",
	Long:  "Reboot or Shutdown device",
	RunE: func(cmd *cobra.Command, args []string) error {
		device := util.GetDeviceByUdId(udid)
		if device == nil {
			os.Exit(0)
		}
		var errReboot error
		if isShutdown {
			errReboot = device.Shutdown()
		} else {
			errReboot = device.Reboot()
		}
		if errReboot != nil {
			fmt.Println("reboot failed")
			os.Exit(0)
		}
		fmt.Println("reboot successful")
		return nil
	},
}

var isShutdown bool

func init() {
	rootCmd.AddCommand(rebootCmd)
	rebootCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber")
	rebootCmd.Flags().BoolVarP(&isShutdown, "shutdown", "s", false, "shutdown your device")
}
