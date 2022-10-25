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
package util

import (
	"archive/zip"
	"bufio"
	"fmt"
	giDevice "github.com/SonicCloudOrg/sonic-gidevice"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const DownLoadTimeOut = 30 * time.Second

var versionMap = map[string]string{
	"12.5": "12.4",
}

var urlList = [...]string{"https://tool.appetizer.io/JinjunHan", "https://code.aliyun.com/hanjinjun", "https://github.com/JinjunHan"}

func GetDeviceByUdId(udId string) (device giDevice.Device) {
	usbMuxClient, err := giDevice.NewUsbmux()
	if err != nil {
		NewErrorPrint(ErrConnect, "usbMux", err)
		return nil
	}
	list, err1 := usbMuxClient.Devices()
	if err1 != nil {
		NewErrorPrint(ErrSendCommand, "listDevices", err1)
		return nil
	}
	if len(list) != 0 {
		if len(udId) != 0 {
			for i, d := range list {
				if d.Properties().SerialNumber == udId {
					device = list[i]
					break
				}
			}
		} else {
			device = list[0]
		}
		if device == nil || device.Properties().SerialNumber == "" {
			fmt.Println("device no found")
			return nil
		}
	} else {
		fmt.Println("no device connected")
		return nil
	}
	return
}

func downloadZip(url, version string) (string, error) {
	vm := version
	if versionMap[version] != "" {
		vm = versionMap[version]
	}
	f, err := os.Stat(".sib")
	if err != nil {
		os.MkdirAll(".sib", os.ModePerm)
		f, err = os.Stat(".sib")
	}
	localAbs, _ := filepath.Abs(f.Name())
	_, errT := os.Stat(fmt.Sprintf(".sib/%s.zip", version))
	if errT != nil {
		client := http.Client{
			Timeout: DownLoadTimeOut,
		}
		res, err := client.Get(fmt.Sprintf("%s/iOSDeviceSupport/raw/master/DeviceSupport/%s.zip", url, vm))
		if err != nil {
			return "", err
		}
		defer res.Body.Close()
		r := bufio.NewReaderSize(res.Body, 32*1024)
		newFile, err := os.Create(fmt.Sprintf(".sib/%s.zip", version))
		w := bufio.NewWriter(newFile)
		io.Copy(w, r)
		abs, _ := filepath.Abs(newFile.Name())
		errZip := unzip(abs, ".sib", version)
		if errZip != nil {
			os.Remove(newFile.Name())
			return "", errZip
		}
	}
	return localAbs, nil
}

func unzip(zipFile, destDir, version string) error {
	zipReader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		var fpath string
		if strings.HasPrefix(f.Name, version) && f.FileInfo().IsDir() {
			fpath = filepath.Join(destDir, version)
		} else {
			fpath = filepath.Join(destDir, version+"/"+path.Base(f.Name))
		}
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}

			inFile, err := f.Open()
			if err != nil {
				return err
			}
			defer inFile.Close()

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func loadDevelopImage(version string) (string, bool) {
	var done = false
	var path = ""
	for _, s := range urlList {
		p, err1 := downloadZip(s, version)
		if err1 == nil {
			path = p
			done = true
			break
		}
	}
	return path, done
}

func GetDeviceVersion(device giDevice.Device) string {
	value, err3 := device.GetValue("", "ProductVersion")
	if err3 != nil {
		NewErrorPrint(ErrSendCommand, "get value", err3)
		os.Exit(0)
	}
	ver := strings.Split(value.(string), ".")
	var reVer string
	if len(ver) >= 2 {
		reVer = ver[0] + "." + ver[1]
	}
	return reVer
}

func CheckMount(device giDevice.Device) {
	sign, errImage := device.Images()
	if errImage != nil || len(sign) == 0 {
		fmt.Println("try to mount developer disk image...")

		reVer := GetDeviceVersion(device)

		p, done := loadDevelopImage(reVer)
		if done {
			var dmg = "DeveloperDiskImage.dmg"
			var sign = dmg + ".signature"
			err4 := device.MountDeveloperDiskImage(fmt.Sprintf("%s/%s/%s", p, reVer, dmg), fmt.Sprintf("%s/%s/%s", p, reVer, sign))
			if err4 != nil {
				fmt.Printf("mount develop disk image fail: %s", err4)
				os.Exit(0)
			}
		} else {
			fmt.Println("download develop disk image fail")
			os.Exit(0)
		}
	}
}

func GetApplicationPID(device giDevice.Device, appName string) (pid int, err error) {
	processes, err := device.AppRunningProcesses()
	if err != nil {
		return -1, err
	}
	for _, p := range processes {
		if p.Name == appName {
			return p.Pid, nil
		}
	}
	return -1, nil
}
