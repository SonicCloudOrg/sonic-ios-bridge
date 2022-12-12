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
package afc

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var afcMkDirCmd = &cobra.Command{
	Use:   "mkdir",
	Short: "create a directory",
	Long:  "create a directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		afcServer := getAFCServer()
		err := (afcServer).Mkdir(mkDir)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		fmt.Println("mkdir success")
		return nil
	},
}

var mkDir string

func initMkDir() {
	afcRootCMD.AddCommand(afcMkDirCmd)
	afcMkDirCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber ( default first device )")
	afcMkDirCmd.Flags().StringVarP(&bundleId, "bundleId", "b", "", "app bundleId")
	afcMkDirCmd.Flags().StringVarP(&mkDir, "folder", "f", "", "mkdir directory path")
	afcMkDirCmd.MarkFlagRequired("folder")
}
