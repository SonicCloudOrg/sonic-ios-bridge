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
	"io"
	"os"
	gPath "path"
)

var afcPushCmd = &cobra.Command{
	Use:   "push",
	Short: "push a file or directory to the device",
	Long:  "push a file or directory to the device",
	RunE: func(cmd *cobra.Command, args []string) error {
		afcServer := getAFCServer()
		pushOperate(afcServer, pushLocalPath, pushSaveDevicePath)
		fmt.Println(fmt.Sprintf("success,push %s --> %s", pushLocalPath, pushSaveDevicePath))

		return nil
	},
}

var pushLocalPath string
var pushSaveDevicePath string

func initPush() {
	afcRootCMD.AddCommand(afcPushCmd)
	afcPushCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber ( default first device )")
	afcPushCmd.Flags().StringVarP(&bundleId, "bundleId", "b", "", "app bundleId")
	afcPushCmd.Flags().StringVarP(&pushLocalPath, "local-path", "l", "", "push file or directory local path")
	afcPushCmd.Flags().StringVarP(&pushSaveDevicePath, "device-path", "d", "", "push save file or directory to device path")
	afcPushCmd.MarkFlagRequired("local-path")
	afcPushCmd.MarkFlagRequired("device-path")
}

func pushFile(afc giDevice.Afc, localPath string, devicePath string) {
	file, err := os.Open(localPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer func() {
		_ = file.Close()
	}()

	afcFile, err := afc.Open(devicePath, giDevice.AfcFileModeWr)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	defer func() {
		_ = afcFile.Close()
	}()
	if _, err = io.Copy(afcFile, file); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

}

func pushOperate(afc giDevice.Afc, localPath string, devicePath string) {
	localFile, err := os.Stat(localPath)
	if err != nil {
		os.Exit(0)
	}
	if localFile.IsDir() {
		aPathInfo, _ := afc.ReadDir(devicePath)
		if aPathInfo == nil {
			mkdirError := afc.Mkdir(devicePath)
			if mkdirError != nil {
				fmt.Println(mkdirError)
				os.Exit(0)
			}
		}
		childFiles, err := os.ReadDir(localPath)
		if err != nil {
			os.Exit(0)
		}
		for _, childFile := range childFiles {
			pushOperate(afc, gPath.Join(localPath, childFile.Name()), gPath.Join(devicePath, childFile.Name()))
		}
	} else {
		pushFile(afc, localPath, devicePath)
	}
}
