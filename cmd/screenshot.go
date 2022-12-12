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
