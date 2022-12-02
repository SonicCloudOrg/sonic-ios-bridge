/*
 *   sonic-ios-bridge  Connect to your iOS Devices.
 *   Copyright (C) 2022 SonicCloudOrg
 *
 *   This program is free software: you can redistribute it and/or modify
 *   it under the terms of the GNU General Public License as published by
 *   the Free Software Foundation, either version 3 of the License, or
 *   (at your option) any later version.
 *
 *   This program is distributed in the hope that it will be useful,
 *   but WITHOUT ANY WARRANTY; without even the implied warranty of
 *   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *   GNU General Public License for more details.
 *
 *   You should have received a copy of the GNU General Public License
 *   along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */
package remote

import (
	"bufio"
	"encoding/json"
	"fmt"
	giDevice "github.com/SonicCloudOrg/sonic-gidevice"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/entity"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
)

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "connect remote device",
	Long:  "connect remote device",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, version, err := CheckRemoteConnect(ip, port, 7)
		if err != nil {
			log.Panic(fmt.Sprintf("connection %s:%d failed,error:%v", ip, port, err))
		}
		log.Printf("connection %s:%d succeeded, device version is:%v", ip, port, version)
		_, err = os.Stat(".sib")
		if err != nil {
			os.MkdirAll(".sib", os.ModePerm)
			os.Stat(".sib")
		}

		file, err := os.OpenFile(RemoteInfoFilePath, os.O_WRONLY|os.O_CREATE, 0666)
		defer file.Close()

		if err != nil {
			log.Panic(err)
		}
		jsonData, err1 := ioutil.ReadFile(RemoteInfoFilePath)
		if err1 != nil {
			log.Panic(err1)
		}

		remoteMap := make(map[string]*entity.RemoteInfo)

		if jsonData != nil && len(jsonData) != 0 {
			err = json.Unmarshal(jsonData, &remoteMap)
			if err != nil {
				log.Panic(err)
			}
		}
		remoteMap[fmt.Sprintf("%s:%d", ip, port)] = &entity.RemoteInfo{
			IP:   &ip,
			Port: &port,
			//Status: OnLine,
		}

		write := bufio.NewWriter(file)

		jsonData, _ = json.Marshal(remoteMap)

		write.Write(jsonData)
		write.Flush()
		return nil
	},
}

func connectInit() {
	remoteCmd.AddCommand(connectCmd)
	connectCmd.Flags().StringVarP(&ip, "ip", "i", "", "remote device ip")
	connectCmd.Flags().IntVarP(&port, "port", "p", 9123, "share port ( default port 9123 )")
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
