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
	giDevice "github.com/SonicCloudOrg/sonic-gidevice"
	"github.com/spf13/cobra"
	"os"
	gPath "path"
)

var afcRMTreeCmd = &cobra.Command{
	Use:   "rmtree",
	Short: "recursively delete all files in a directory",
	Long:  "recursively delete all files in a directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		afcServer := getAFCServer()
		removeTree(afcServer, rmDir)
		fmt.Println("success")
		return nil
	},
}

var rmDir string

func initRMTree() {
	afcRootCMD.AddCommand(afcRMTreeCmd)
	afcRMTreeCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber ( default first device )")
	afcRMTreeCmd.Flags().StringVarP(&bundleId, "bundleId", "b", "", "app bundleId")
	afcRMTreeCmd.Flags().StringVarP(&rmDir, "folder", "f", "", "folder address to delete")
	afcRMTreeCmd.MarkFlagRequired("folder")
}

func removeTree(afc giDevice.Afc, devicePath string) {
	fileInfo, err := afc.Stat(devicePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	if fileInfo.IsDir() {
		fileNames, err := afc.ReadDir(devicePath)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
		for _, file := range fileNames {
			if file == "." || file == ".." {
				continue
			}
			var childPath string
			childPath = gPath.Join(devicePath, file)

			removeTree(afc, childPath)
		}

		err = afc.Remove(devicePath)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	} else {
		err := afc.Remove(devicePath)
		if err != nil {
			fmt.Println(err)
			os.Exit(0)
		}
	}
}
