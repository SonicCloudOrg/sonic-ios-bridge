package conn

import (
	"bytes"
	"encoding/json"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/tool"
	"howett.net/plist"
	"strings"
)

type DeviceList struct {
	DeviceList []iDevice `json:"deviceList"`
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

func (usbMuxClient *UsbMuxClient) ListDevices() (DeviceList, error) {
	err := usbMuxClient.Send(NewListDevicesMessage())
	if err != nil {
		return DeviceList{}, tool.NewErrorPrint(tool.ErrSendCommand, "listDevices", err)
	}
	resp, err := usbMuxClient.ReadMessage()
	if err != nil {
		return DeviceList{}, tool.NewErrorPrint(tool.ErrReadingMsg, "deviceList", err)
	}
	return deviceListForBytes(resp.Payload), nil
}

func deviceListForBytes(plistBytes []byte) DeviceList {
	decoder := plist.NewDecoder(bytes.NewReader(plistBytes))
	var deviceList DeviceList
	decoder.Decode(&deviceList)
	for i, d := range deviceList.DeviceList {
		deviceList.DeviceList[i].Status = d.GetStatus()
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

func (deviceList DeviceList) ToString() string {
	var s strings.Builder
	for _, e := range deviceList.DeviceList {
		s.WriteString(e.Properties.SerialNumber + " " + e.Status)
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
