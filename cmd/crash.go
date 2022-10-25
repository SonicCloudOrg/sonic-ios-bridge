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
	giDevice "github.com/SonicCloudOrg/sonic-gidevice"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var crashCmd = &cobra.Command{
	Use:   "crash",
	Short: "Get CrashReport from your device",
	Long:  "Get CrashReport from your device",
	RunE: func(cmd *cobra.Command, args []string) error {
		device := util.GetDeviceByUdId(udid)
		if device == nil {
			os.Exit(0)
		}
		if !filepath.IsAbs(crashOutputPath) {
			var err error
			if crashOutputPath, err = filepath.Abs(crashOutputPath); err != nil {
				fmt.Println("path no found!")
				os.Exit(0)
			}
		}
		err := device.MoveCrashReport(crashOutputPath,
			giDevice.WithKeepCrashReport(keep),
			giDevice.WithExtractRawCrashReport(true),
			giDevice.WithWhenMoveIsDone(func(filename string) {
				fmt.Printf("%s: done.\n", filename)
			}),
		)
		if err != nil {
			return util.NewErrorPrint(util.ErrUnknown, "move crash files", err)
		}
		fmt.Println("All done.")
		return nil
	},
}

var keep bool
var crashOutputPath string

func init() {
	rootCmd.AddCommand(crashCmd)
	crashCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber (default first device)")
	crashCmd.Flags().BoolVarP(&keep, "keep", "k", false, "keep crash reports from device")
	crashCmd.Flags().StringVarP(&crashOutputPath, "path", "p", "./", "output path")
}
