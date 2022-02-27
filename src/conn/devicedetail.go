package conn

import (
	"encoding/json"
	"fmt"
	"github.com/Masterminds/semver"
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

type GetValueRequest struct {
	Label   string
	Key     string `plist:"Key,omitempty"`
	Request string `plist:"Request"`
	Domain  string `plist:"Domain,omitempty"`
	Value   string `plist:"Value,omitempty"`
}

func NewGetValue(domain string, key string) GetValueRequest {
	data := GetValueRequest{}

	if len(domain) > 0 {
		data = GetValueRequest{
			Label:   BundleId,
			Domain:  domain,
			Key:     key,
			Request: "GetValue",
		}
	} else {
		data = GetValueRequest{
			Label:   BundleId,
			Key:     key,
			Request: "GetValue",
		}
	}
	return data
}

func GetValueFromDevice(device iDevice, domain, key string) (interface{}, error) {
	lockdownConnection, err := NewLockdownConnection(device)
	if err != nil {
		return map[string]interface{}{}, err
	}
	defer lockdownConnection.Close()
	err = lockdownConnection.Send(NewGetValue(domain, key))
	if err != nil {
		return map[string]interface{}{}, err
	}
	resp, err := lockdownConnection.ReadMessage()
	if err != nil {
		return map[string]interface{}{}, err
	}
	plist, err := parsePlist(resp)
	if err != nil {
		return map[string]interface{}{}, err
	}
	return plist["Value"], err
}

func (device *iDevice) GetDetail() (*DeviceDetail, error) {
	values, err := GetValueFromDevice(*device, "", "")
	if err != nil {
		return &DeviceDetail{}, err
	}
	data, _ := json.Marshal(values)
	detail := &DeviceDetail{}
	json.Unmarshal(data, detail)
	detail.GenerationName = detail.GetGenerationName()
	return detail, nil
}

func (device *iDevice) GetProductVersion() (string, error) {
	values, err := GetValueFromDevice(*device, "", "ProductVersion")
	if err != nil {
		return "", err
	}
	return values.(string), nil
}

func (device *iDevice) GetSemverProductVersion() (*semver.Version, error) {
	values, err := device.GetProductVersion()
	if err != nil {
		return &semver.Version{}, err
	}
	version, err := semver.NewVersion(values)
	return version, err
}

//

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
