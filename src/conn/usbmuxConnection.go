package conn

import "fmt"

type UsbMuxConnection struct {
	deviceConnection DeviceConnectInterface
	tag              uint32
}

func NewUsbMuxConnection(deviceConnection DeviceConnectInterface) *UsbMuxConnection {
	usbMuxConn := &UsbMuxConnection{tag: 0, deviceConnection: deviceConnection}
	var err error
	usbMuxConn.deviceConnection,err = NewDeviceConnection()
	if err != nil {
		fmt.Errorf("fail to get device connection: %w", err)
		return nil
	}
	return usbMuxConn
}
