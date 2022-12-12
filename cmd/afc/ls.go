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

var afcLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "ls to view the directory",
	Long:  "ls to view the directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		afcServer := getAFCServer()
		lsShow(afcServer, lsDirPath)
		return nil
	},
}

var lsDirPath string

func initLs() {
	afcRootCMD.AddCommand(afcLsCmd)
	afcLsCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber ( default first device )")
	afcLsCmd.Flags().StringVarP(&bundleId, "bundleId", "b", "", "app bundleId")
	afcLsCmd.Flags().StringVarP(&lsDirPath, "folder", "f", "", "ls folder path")
	afcLsCmd.MarkFlagRequired("folder")
}

func lsShow(afc giDevice.Afc, filePath string) {
	fileNames, err := afc.ReadDir(filePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	for _, fileName := range fileNames {
		if fileName == "." || fileName == ".." {
			continue
		}
		info, err := afc.Stat(gPath.Join(filePath, fileName))

		if err != nil {
			os.Exit(0)
		}
		if info.IsDir() {
			fmt.Println(fileName + "/")
		} else {
			fmt.Println(fmt.Sprintf("- %s %d", fileName, info.Size()))
		}
	}
}
