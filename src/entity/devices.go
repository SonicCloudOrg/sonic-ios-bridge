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
