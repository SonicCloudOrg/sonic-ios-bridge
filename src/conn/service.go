package conn

import (
	"bytes"
	"fmt"
	"howett.net/plist"
)

type StartServiceRequest struct {
	Request string `plist:"Request"`
	Label   string
	Service string
}

type StartServiceResp struct {
	Request          string `plist:"Request"`
	Port             uint16
	Service          string
	EnableServiceSSL bool
	Error            string
}

func NewStartServiceRequest(serviceName string) StartServiceRequest {
	return StartServiceRequest{
		Request: "StartService",
		Label:   BundleId,
		Service: serviceName,
	}
}

func startServiceRespForBytes(plistBytes []byte) StartServiceResp {
	decoder := plist.NewDecoder(bytes.NewReader(plistBytes))
	var data StartServiceResp
	_ = decoder.Decode(&data)
	return data
}

func StartServiceFromDevice(device iDevice, serviceName string) (StartServiceResp, error) {
	lockdownConnection, err := NewLockdownConnection(device)
	if err != nil {
		return StartServiceResp{}, err
	}
	defer lockdownConnection.Close()
	err = lockdownConnection.Send(NewStartServiceRequest(serviceName))
	if err != nil {
		return StartServiceResp{}, err
	}
	resp, err := lockdownConnection.ReadMessage()
	if err != nil {
		return StartServiceResp{}, err
	}
	response := startServiceRespForBytes(resp)
	if response.Error != "" {
		return StartServiceResp{}, fmt.Errorf("could not start service %s", serviceName)
	}
	fmt.Println("device started service")
	return response, nil
}

func (device *iDevice) StartService(serviceName string) (*StartServiceResp, error) {
	response, err := StartServiceFromDevice(*device, serviceName)
	if err != nil {
		return &response, err
	}
	return &response, nil
}
