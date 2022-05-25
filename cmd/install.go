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
	"github.com/spf13/cobra"
	"os"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install App in your device",
	Long:  "Install App in your device",
	RunE: func(cmd *cobra.Command, args []string) error {
		device := util.GetDeviceByUdId(udid)
		if device == nil {
			os.Exit(0)
		}
		errInstall := device.AppInstall(path)
		if errInstall != nil {
			fmt.Printf("install failed: %s", errInstall)
			os.Exit(0)
		}
		fmt.Println("install successful")
		return nil
	},
}

var path string

func init() {
	appCmd.AddCommand(installCmd)
	installCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber")
	installCmd.Flags().StringVarP(&path, "path", "p", "", "path of ipa file")
	installCmd.MarkFlagRequired("path")
}
