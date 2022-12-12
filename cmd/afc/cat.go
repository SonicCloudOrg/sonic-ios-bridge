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
)

var afcCatCmd = &cobra.Command{
	Use:   "cat",
	Short: "cat to view files",
	Long:  "cat to view files",
	RunE: func(cmd *cobra.Command, args []string) error {
		afcServer := getAFCServer()
		catFile(afcServer, catFilePath)
		return nil
	},
}

var catFilePath string

func initCat() {
	afcRootCMD.AddCommand(afcCatCmd)
	afcCatCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber ( default first device )")
	afcCatCmd.Flags().StringVarP(&bundleId, "bundleId", "b", "", "app bundleId")
	afcCatCmd.Flags().StringVarP(&catFilePath, "file", "f", "", "cat file path")
	afcCatCmd.MarkFlagRequired("file")
}

func catFile(afc giDevice.Afc, filePath string) {
	fileInfo, err := afc.Stat(filePath)
	if err != nil {
		fmt.Println("file path is null")
		os.Exit(0)
	}
	p := make([]byte, fileInfo.Size())
	afcFile, err := afc.Open(filePath, giDevice.AfcFileModeRdOnly)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	defer func() {
		_ = afcFile.Close()
	}()
	n, err := afcFile.Read(p)
	if err == io.EOF {
		fmt.Println(err)
		os.Exit(0)
		//break
	}
	fmt.Print(string(p[:n]))
}
