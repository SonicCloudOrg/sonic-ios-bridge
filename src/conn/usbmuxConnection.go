package conn

type UsbMuxConnection struct {
	deviceConnection DeviceConnectInterface
	tag              uint32
}

func NewUsbMuxConnection(deviceConnection DeviceConnectInterface) *UsbMuxConnection {
	usbMuxConn := &UsbMuxConnection{tag: 0, deviceConnection: deviceConnection}
	return usbMuxConn
}
