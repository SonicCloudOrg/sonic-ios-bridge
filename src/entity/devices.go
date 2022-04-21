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
package entity

import (
	"encoding/json"
	"strings"
)

type DeviceList struct {
	DeviceList []Device `json:"deviceList"`
}

type Device struct {
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
	var s strings.Builder
	for _, e := range deviceList.DeviceList {
		s.WriteString(e.SerialNumber + " " + e.Status)
	}
	return s.String()
}

func (deviceList DeviceList) ToJson() string {
	result, _ := json.Marshal(deviceList)
	return string(result)
}

func (deviceList DeviceList) ToFormat() string {
	result, _ := json.MarshalIndent(deviceList, "", "\t")
	return string(result)
}
