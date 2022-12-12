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
package entity

import (
	"encoding/json"
	"fmt"
	"strings"
)

type DeviceList struct {
	DeviceList []Device `json:"deviceList"`
}

type Device struct {
	RemoteAddr      string       `json:"remoteAddr"`
	DeviceID        int          `json:"deviceId"`
	ConnectionSpeed int          `json:"connectionSpeed"`
	ConnectionType  string       `json:"connectionType"`
	LocationID      int          `json:"locationId"`
	ProductID       int          `json:"productId"`
	SerialNumber    string       `json:"serialNumber"`
	Status          string       `json:"status"`
	DeviceDetail    DeviceDetail `json:"deviceDetail"`
}

func (device *Device) GetStatus() string {
	if device.ConnectionType != "" {
		return "online"
	} else {
		return "offline"
	}
}

func (device Device) ToString() string {
	var s strings.Builder
	s.WriteString(device.SerialNumber + " " + device.Status)
	return s.String()
}

func (device Device) ToJson() string {
	result, _ := json.Marshal(device)
	return string(result)
}

func (device Device) ToFormat() string {
	result, _ := json.MarshalIndent(device, "", "\t")
	return string(result)
}

func (deviceList DeviceList) ToString() string {
	for _, e := range deviceList.DeviceList {
		fmt.Println(e.SerialNumber + " " + e.Status + " " + e.RemoteAddr)
	}
	return ""
}

func (deviceList DeviceList) ToJson() string {
	result, _ := json.Marshal(deviceList)
	return string(result)
}

func (deviceList DeviceList) ToFormat() string {
	result, _ := json.MarshalIndent(deviceList, "", "\t")
	return string(result)
}
