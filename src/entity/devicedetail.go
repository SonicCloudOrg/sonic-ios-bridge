package entity

import (
	"encoding/json"
	"fmt"
	giDevice "github.com/electricbubble/gidevice"
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
