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
package util

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	giDevice "github.com/SonicCloudOrg/sonic-gidevice"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/entity"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	DownLoadTimeOut    = 30 * time.Second
	baseDir            = ".sib"
	RemoteInfoFilePath = baseDir + string(filepath.Separator) + "connect.txt"
)

var versionMap = map[string]string{
	"12.5": "12.4",
}

var urlList = [...]string{"https://tool.appetizer.io/JinjunHan", "https://code.aliyun.com/hanjinjun", "https://github.com/JinjunHan"}

func GetDeviceByUdId(udId string) (device giDevice.Device) {
	remoteList, err2 := ReadRemote()
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
	if len(list) != 0 || len(remoteList) != 0 {
		if len(list) == 0 {
			list = []giDevice.Device{}
		}
		if err2 == nil {
			for _, v := range remoteList {
				list = append(list, v)
			}
		}
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

func ReadRemote() (remoteDevList map[string]giDevice.Device, err error) {
	defer func() {

		if r := recover(); r != nil {
			fmt.Println("recover...:", r)
		}
	}()

	file, err := os.Open(RemoteInfoFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	if content == nil && len(content) != 0 {
		return nil, errors.New("remote info file non existent remote data")
	}
	remoteInfoData := make(map[string]entity.RemoteInfo)
	err = json.Unmarshal(content, &remoteInfoData)
	if err != nil {
		//fmt.Println(err)
		return nil, err
	}

	if remoteDevList == nil {
		remoteDevList = map[string]giDevice.Device{}
	}
	var wait sync.WaitGroup
	for k, v := range remoteInfoData {
		//if v.Status!=remote.OnLine {
		//	continue
		//}
		wait.Add(1)
		go func(info entity.RemoteInfo) {
			dev, _, err1 := CheckRemoteConnect(*info.Host, *info.Port, 5)
			if err1 != nil {
				wait.Done()
				return
			}
			remoteDevList[k] = dev
			wait.Done()
		}(v)
	}
	wait.Wait()
	return remoteDevList, nil
}

func CheckRemoteConnect(ip string, port int, timeout int) (dev giDevice.Device, version interface{}, err error) {
	dev, err = giDevice.NewRemoteConnect(ip, port, timeout)
	if err != nil {
		return nil, nil, err
	}
	version, err = dev.GetValue("", "ProductVersion")
	if err != nil {
		return nil, nil, err
	}
	return dev, version, nil
}

func downloadZip(url, version string) (string, error) {
	vm := version
	if versionMap[version] != "" {
		vm = versionMap[version]
	}
	f, err := os.Stat(baseDir)
	if err != nil {
		os.MkdirAll(baseDir, os.ModePerm)
		f, err = os.Stat(baseDir)
	}
	localAbs, _ := filepath.Abs(f.Name())
	filePath := fmt.Sprintf("%s.zip", baseDir+string(filepath.Separator)+version)
	_, errT := os.Stat(filePath)
	if errT != nil {
		client := http.Client{
			Timeout: DownLoadTimeOut,
		}
		res, err := client.Get(fmt.Sprintf("%s/iOSDeviceSupport/raw/master/iOSDeviceSupport/%s.zip", url, vm))
		if err != nil {
			return "", err
		}
		defer res.Body.Close()
		r := bufio.NewReaderSize(res.Body, 32*1024)
		newFile, err := os.Create(filePath)
		w := bufio.NewWriter(newFile)
		io.Copy(w, r)
		abs, _ := filepath.Abs(newFile.Name())
		errZip := unzip(abs, baseDir, version)
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
			fpath = filepath.Join(destDir, version+string(filepath.Separator)+path.Base(f.Name))
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

func StartProxy() func(listener net.Listener, port int, device giDevice.Device) {
	return func(listener net.Listener, port int, device giDevice.Device) {
		for {
			var accept net.Conn
			var err error
			if accept, err = listener.Accept(); err != nil {
				log.Println("accept:", err)
			}
			fmt.Println("accept", accept.RemoteAddr())
			rInnerConn, err := device.NewConnect(port)
			var retry = 0
			for {
				retry++
				if retry > 5 {
					break
				}
				if err != nil {
					fmt.Println("connect to device fail...retry in 2s...")
					time.Sleep(time.Duration(2) * time.Second)
					rInnerConn, err = device.NewConnect(port)
				} else {
					break
				}
			}
			rConn := rInnerConn.RawConn()
			rConn.SetDeadline(time.Time{})
			go func(lConn net.Conn) {
				go func(lConn, rConn net.Conn) {
					if _, err := io.Copy(lConn, rConn); err != nil {
						log.Println("local -> remote failed:", err)
					}
				}(lConn, rConn)
				go func(lConn, rConn net.Conn) {
					if _, err := io.Copy(rConn, lConn); err != nil {
						log.Println("local <- remote failed:", err)
					}
				}(lConn, rConn)
			}(accept)
		}
	}
}
