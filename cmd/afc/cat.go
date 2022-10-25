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
