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
	"fmt"
	giDevice "github.com/SonicCloudOrg/sonic-gidevice"
)

type DeviceDetail struct {
	GenerationName            string `json:"generationName,omitempty"`
	DeviceName                string `json:"deviceName,omitempty"`
	DeviceColor               string `json:"deviceColor,omitempty"`
	DeviceClass               string `json:"deviceClass,omitempty"`
	ProductVersion            string `json:"productVersion,omitempty"`
	ProductType               string `json:"productType,omitempty"`
	ProductName               string `json:"productName,omitempty"`
	PasswordProtected         bool   `json:"passwordProtected,omitempty"`
	ModelNumber               string `json:"modelNumber,omitempty"`
	SerialNumber              string `json:"serialNumber,omitempty"`
	SIMStatus                 string `json:"simStatus,omitempty"`
	PhoneNumber               string `json:"phoneNumber,omitempty"`
	CPUArchitecture           string `json:"cpuArchitecture,omitempty"`
	ProtocolVersion           string `json:"protocolVersion,omitempty"`
	RegionInfo                string `json:"regionInfo,omitempty"`
	TelephonyCapability       bool   `json:"telephonyCapability,omitempty"`
	TimeZone                  string `json:"timeZone,omitempty"`
	UniqueDeviceID            string `json:"uniqueDeviceID,omitempty"`
	WiFiAddress               string `json:"wifiAddress,omitempty"`
	WirelessBoardSerialNumber string `json:"wirelessBoardSerialNumber,omitempty"`
	BluetoothAddress          string `json:"bluetoothAddress,omitempty"`
	BuildVersion              string `json:"buildVersion,omitempty"`
}

func GetDetail(device giDevice.Device) (*DeviceDetail, error) {
	value, err1 := device.GetValue("", "")
	if err1 != nil {
		return &DeviceDetail{}, fmt.Errorf("get %s device detail fail : %w", device.Properties().SerialNumber, err1)
	}
	detailByte, _ := json.Marshal(value)
	detail := &DeviceDetail{}
	json.Unmarshal(detailByte, detail)
	detail.GenerationName = detail.GetGenerationName()
	return detail, nil
}
