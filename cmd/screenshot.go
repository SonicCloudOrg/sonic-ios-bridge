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
	uuid "github.com/satori/go.uuid"
	"github.com/spf13/cobra"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
)

var screenshotCmd = &cobra.Command{
	Use:   "screenshot",
	Short: "Get screenshot realtime",
	Long:  "Get screenshot realtime",
	RunE: func(cmd *cobra.Command, args []string) error {
		device := util.GetDeviceByUdId(udid)
		if device == nil {
			os.Exit(0)
		}
		util.CheckMount(device)
		bytes, err := device.Screenshot()
		if err != nil {
			return util.NewErrorPrint(util.ErrSendCommand, "screenshot", err)
		}
		img, _, errDecode := image.Decode(bytes)
		if errDecode != nil {
			return util.NewErrorPrint(util.ErrUnknown, "decode img", err)
		}
		if len(fileName) == 0 {
			fileName = uuid.NewV4().String()
		}
		outputFile := filePath + "/" + fileName + "." + fileType
		file, err := os.Create(outputFile)
		if err != nil {
			return util.NewErrorPrint(util.ErrUnknown, "create file", err)
		}
		defer file.Close()
		switch fileType {
		case "png":
			err = png.Encode(file, img)
		case "jpeg":
			err = jpeg.Encode(file, img, nil)
		default:
			err = png.Encode(file, img)
		}
		if err != nil {
			return util.NewErrorPrint(util.ErrUnknown, "encode img", err)
		}
		fileInfo, _ := file.Stat()
		pathResult, _ := filepath.Abs(filePath)
		fmt.Printf("screenshot to path:%s name:%s size:%dbytes\n", pathResult, fileInfo.Name(), fileInfo.Size())
		return nil
	},
}

var filePath, fileName, fileType string

func init() {
	rootCmd.AddCommand(screenshotCmd)
	screenshotCmd.Flags().StringVarP(&udid, "udid", "u", "", "device's serialNumber ( default first device )")
	screenshotCmd.Flags().StringVarP(&filePath, "path", "p", "./", "output path")
	screenshotCmd.Flags().StringVarP(&fileName, "name", "n", "", "output file name")
	screenshotCmd.Flags().StringVarP(&fileType, "type", "t", "png", "output file format (png or jpeg)")
}
