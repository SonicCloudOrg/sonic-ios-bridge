package conn

import (
	"fmt"
	"net"
)

const LockdownPort uint16 = 62078

type LockDownConnection struct {
	sessionID        string
	plistCodec       PlistCodec
	deviceConnection DeviceConnectInterface
}

type ConnectMessage struct {
	BundleID            string
	ClientVersionString string
	MessageType         string
	ProgName            string
	DeviceID            uint32
	PortNumber          uint16
	LibUSBMuxVersion    uint32 `plist:"kLibUSBMuxVersion"`
}

func NewConnectMessage(deviceID int, portNumber uint16) ConnectMessage {
	data := ConnectMessage{
		MessageType:         "Connect",
		BundleID:            BundleId,
		ClientVersionString: ClientVersion,
		ProgName:            ProgramName,
		DeviceID:            uint32(deviceID),
		PortNumber:          portNumber,
		LibUSBMuxVersion:    3,
	}
	return data
}

func NewLockdownConnection(device iDevice) (*LockDownConnection, error) {
	usbMuxClient, err := NewUsbMuxClient()
	if err != nil {
		return nil, err
	}
	pairRecord, err := usbMuxClient.ReadPair(device.Properties.SerialNumber)
	if err != nil {
		return nil, err
	}
	lockdownConnection, err := usbMuxClient.ConnectLockdown(device.DeviceID)
	if err != nil {
		return nil, err
	}
	_, err = lockdownConnection.StartSession(pairRecord)
	if err != nil {
		return nil, fmt.Errorf("startSession fail: %w", err)
	}
	return lockdownConnection, nil
}

func (usbMuxClient *UsbMuxClient) ConnectLockdown(deviceID int) (*LockDownConnection, error) {
	msg := NewConnectMessage(deviceID, LockdownPort)
	err := usbMuxClient.Send(msg)
	if err != nil {
		return &LockDownConnection{}, err
	}
	resp, err := usbMuxClient.ReadMessage()
	if err != nil {
		return &LockDownConnection{}, err
	}
	if !usbMuxRespForBytes(resp.Payload).IsSuccess() {
		return &LockDownConnection{"", NewPlistCodec(), usbMuxClient.deviceConnection}, fmt.Errorf("fail connect to lockdown")
	}
	return nil, fmt.Errorf("fail connect to lockdown")
}

func NewLockDownConnection(deviceConnect DeviceConnectInterface) *LockDownConnection {
	return &LockDownConnection{deviceConnection: deviceConnect, plistCodec: NewPlistCodec()}
}

func (lockDownConn *LockDownConnection) Close() {
	lockDownConn.StopSession()
	lockDownConn.deviceConnection.Close()
}

func (lockDownConn LockDownConnection) DisableSessionSSL() {
	lockDownConn.deviceConnection.DisableSessionSSL()
}

func (lockDownConn LockDownConnection) EnableSessionSsl(pairRecord PairRecord) error {
	return lockDownConn.deviceConnection.EnableSessionSSL(pairRecord)
}

func (lockDownConn LockDownConnection) Send(msg interface{}) error {
	bytes, err := lockDownConn.plistCodec.Encode(msg)
	if err != nil {
		return err
	}
	return lockDownConn.deviceConnection.Send(bytes)
}

func (lockDownConn *LockDownConnection) ReadMessage() ([]byte, error) {
	reader := lockDownConn.deviceConnection.Reader()
	resp, err := lockDownConn.plistCodec.Decode(reader)
	if err != nil {
		return make([]byte, 0), err
	}
	return resp, err
}

func (lockDownConn *LockDownConnection) GetConn() net.Conn {
	return lockDownConn.deviceConnection.GetConn()
}

func GetValueFromDevice(device iDevice) (map[string]interface{}, error) {
	lockdownConnection, err := NewLockdownConnection(device)
	if err != nil {
		return map[string]interface{}{}, err
	}
	defer lockdownConnection.Close()
	err = lockdownConnection.Send(NewGetValue("", "ProductVersion"))
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
	plist, ok := plist["Value"].(map[string]interface{})
	if !ok {
		return plist, err
	}
	return plist, err
}
