package conn

import (
	"bytes"
	"encoding/json"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/tool"
	"howett.net/plist"
	"strings"
)

type IDeviceList struct {
	IDeviceList []iDevice `json:"deviceList"`
}

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

type iDevice struct {
	DeviceID     int          `json:"deviceId"`
	MessageType  string       `json:"messageType"`
	Properties   DeviceProp   `json:"properties"`
	Status       string       `json:"status"`
	DeviceDetail DeviceDetail `json:"deviceDetail"`
}

type DeviceProp struct {
	ConnectionSpeed int    `json:"connectionSpeed"`
	ConnectionType  string `json:"connectionType"`
	DeviceID        int    `json:"deviceId"`
	LocationID      int    `json:"locationId"`
	ProductID       int    `json:"productId"`
	SerialNumber    string `json:"serialNumber"`
}

type ListDevicesMessage struct {
	MessageType         string
	ProgName            string
	ClientVersionString string
}

func NewListDevicesMessage() ListDevicesMessage {
	msg := ListDevicesMessage{
		MessageType:         "ListDevices",
		ProgName:            ProgramName,
		ClientVersionString: ClientVersion,
	}
	return msg
}

func (device *iDevice) GetStatus() string {
	if device.MessageType == "Attached" {
		return "online"
	} else {
		return "offline"
	}
}

func (usbMuxClient *UsbMuxClient) ListDevices() (IDeviceList, error) {
	err := usbMuxClient.Send(NewListDevicesMessage())
	if err != nil {
		return IDeviceList{}, tool.NewErrorPrint(tool.ErrSendCommand, "listDevices", err)
	}
	defer usbMuxClient.GetDeviceConn().Close()
	resp, err := usbMuxClient.ReadMessage()
	if err != nil {
		return IDeviceList{}, tool.NewErrorPrint(tool.ErrReadingMsg, "deviceList", err)
	}
	return deviceListForBytes(resp.Payload), nil
}

func deviceListForBytes(plistBytes []byte) IDeviceList {
	decoder := plist.NewDecoder(bytes.NewReader(plistBytes))
	var deviceList IDeviceList
	decoder.Decode(&deviceList)
	for i, d := range deviceList.IDeviceList {
		deviceList.IDeviceList[i].Status = d.GetStatus()
	}
	return deviceList
}

var deviceIdMap = make(map[int]string)

func deviceForBytes(plistBytes []byte) iDevice {
	decoder := plist.NewDecoder(bytes.NewReader(plistBytes))
	var device iDevice
	decoder.Decode(&device)
	device.Status = device.GetStatus()
	if len(device.Properties.SerialNumber) > 0 {
		deviceIdMap[device.DeviceID] = device.Properties.SerialNumber
	} else {
		device.Properties.SerialNumber = deviceIdMap[device.DeviceID]
		delete(deviceIdMap, device.DeviceID)
	}
	return device
}

func (device iDevice) ToString() string {
	var s strings.Builder
	s.WriteString(device.Properties.SerialNumber + " " + device.Status)
	return s.String()
}

func (device iDevice) ToJson() string {
	result, _ := json.Marshal(device)
	return string(result)
}

func (device iDevice) ToFormat() string {
	result, _ := json.MarshalIndent(device, "", "\t")
	return string(result)
}

func (deviceList IDeviceList) ToString() string {
	var s strings.Builder
	for _, e := range deviceList.IDeviceList {
		s.WriteString(e.Properties.SerialNumber + " " + e.Status)
	}
	return s.String()
}

func (deviceList IDeviceList) ToJson() string {
	result, _ := json.Marshal(deviceList)
	return string(result)
}

func (deviceList IDeviceList) ToFormat() string {
	result, _ := json.MarshalIndent(deviceList, "", "\t")
	return string(result)
}

//
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
