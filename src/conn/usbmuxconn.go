package conn

import (
	"fmt"
)

type UsbMuxClient struct {
	deviceConnection DeviceConnectInterface
	tag              uint32
}

func NewUsbMuxClient() (u *UsbMuxClient, err error) {
	deviceConnection, err := NewDeviceConnection()
	if err != nil {
		return nil, fmt.Errorf("fail to get device connection: %w", err)
	}
	u = &UsbMuxClient{tag: 0, deviceConnection: deviceConnection}
	return
}
