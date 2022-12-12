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
