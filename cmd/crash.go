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
