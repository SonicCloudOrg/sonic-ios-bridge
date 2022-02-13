package conn

import (
	"bytes"
	"fmt"
	"howett.net/plist"
	"strings"
)

type DeviceList struct {
	DeviceList []iDevice
}

type iDevice struct {
	DeviceID    int
	MessageType string
	Properties  DeviceProp
}

type DeviceProp struct {
	ConnectionSpeed int
	ConnectionType  string
	DeviceID        int
	LocationID      int
	ProductID       int
	SerialNumber    string
}

type ListDevicesMessage struct {
	MessageType         string
	ProgName            string
	ClientVersionString string
}

func NewListDevicesMessage() ListDevicesMessage {
	data := ListDevicesMessage{
		MessageType:         "ListDevices",
		ProgName:            ProgramName,
		ClientVersionString: ClientVersion,
	}
	return data
}

func (usbMuxClient *UsbMuxClient) ListDevices() (DeviceList, error) {
	err := usbMuxClient.Send(NewListDevicesMessage())
	if err != nil {
		return DeviceList{}, fmt.Errorf("Failed sending to usbmux requesting devicelist: %v", err)
	}
	response, err := usbMuxClient.ReadMessage()
	if err != nil {
		return DeviceList{}, fmt.Errorf("Failed getting devicelist: %v", err)
	}
	return DeviceListForBytes(response.Payload), nil
}

func DeviceListForBytes(plistBytes []byte) DeviceList {
	decoder := plist.NewDecoder(bytes.NewReader(plistBytes))
	var deviceList DeviceList
	decoder.Decode(&deviceList)
	return deviceList
}

func (deviceList DeviceList) ToString() string {
	var s strings.Builder
	for _, e := range deviceList.DeviceList {
		s.WriteString(e.Properties.SerialNumber)
		s.WriteString("\n")
	}
	return s.String()
}

func (deviceList DeviceList) ToJson() map[string]interface{} {
	devices := make([]string, len(deviceList.DeviceList))
	for i, e := range deviceList.DeviceList {
		devices[i] = e.Properties.SerialNumber
	}
	return map[string]interface{}{"list": devices}
}
