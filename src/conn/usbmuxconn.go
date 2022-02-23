package conn

import (
	"encoding/binary"
	"fmt"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/tool"
	"io"
)

const (
	BundleId      = "sib.conn"
	ProgramName   = "sonic-ios-bridge"
	ClientVersion = "sonic-ios-bridge-1.0.0"
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

type UsbMuxMessage struct {
	Header  UsbMuxHeader
	Payload []byte
}

// UsbMuxHeader header for plist
type UsbMuxHeader struct {
	Length  uint32
	Version uint32
	Request uint32
	Tag     uint32
}

func (usbMuxClient *UsbMuxClient) GetDeviceConn() DeviceConnectInterface {
	return usbMuxClient.deviceConnection
}

func (usbMuxClient *UsbMuxClient) Send(msg interface{}) error {
	if usbMuxClient.deviceConnection == nil {
		return io.EOF
	}
	writer := usbMuxClient.deviceConnection.Writer()
	//tag message and prepare to receive
	usbMuxClient.tag++
	err := usbMuxClient.encode(msg, writer)
	if err != nil {
		return fmt.Errorf("error send to usbmux:%w", err)
	}
	return nil
}

func (usbMuxClient *UsbMuxClient) ReadMessage() (UsbMuxMessage, error) {
	if usbMuxClient.deviceConnection == nil {
		return UsbMuxMessage{}, io.EOF
	}
	reader := usbMuxClient.deviceConnection.Reader()
	msg, err := usbMuxClient.decode(reader)
	if err != nil {
		return UsbMuxMessage{}, err
	}
	return msg, nil
}

//https://github.com/danielpaulus/go-ios
func (usbMuxClient *UsbMuxClient) encode(message interface{}, writer io.Writer) error {
	bytes := transToPlistBytes(message)
	err := writeHeader(len(bytes), usbMuxClient.tag, writer)
	if err != nil {
		return err
	}
	_, err = writer.Write(bytes)
	return err
}

func writeHeader(length int, tag uint32, writer io.Writer) error {
	//This tag is which to receive
	header := UsbMuxHeader{Length: 16 + uint32(length), Request: 8, Version: 1, Tag: tag}
	return binary.Write(writer, binary.LittleEndian, header)
}

//https://github.com/alibaba/taobao-iphone-device/blob/main/tidevice/_usbmux.py
func (usbMuxClient *UsbMuxClient) decode(r io.Reader) (UsbMuxMessage, error) {
	var usbMuxHeader UsbMuxHeader
	err := binary.Read(r, binary.LittleEndian, &usbMuxHeader)
	if err != nil {
		return UsbMuxMessage{}, err
	}
	payLoadBytes := make([]byte, usbMuxHeader.Length-16)
	n, err := io.ReadFull(r, payLoadBytes)
	if err != nil {
		return UsbMuxMessage{}, fmt.Errorf("error decode msg %d : %w", n, err)
	}
	return UsbMuxMessage{usbMuxHeader, payLoadBytes}, nil
}

func (usbMuxClient *UsbMuxClient) Connect(deviceId int, port int) error {
	msg := NewConnectMessage(deviceId, port)
	usbMuxClient.Send(msg)
	resp, err := usbMuxClient.ReadMessage()
	if err != nil {
		return err
	}
	response := usbMuxRespForBytes(resp.Payload)
	if response.IsSuccess() {
		return nil
	}
	return tool.NewErrorPrint(tool.ErrConnect, "service", nil)
}

func (usbMuxClient *UsbMuxClient) ConnectStartService(deviceID int, resp StartServiceResp, pairRecord PairRecord) error {
	err := usbMuxClient.Connect(deviceID, int(resp.Port))
	if err != nil {
		return err
	}
	if resp.EnableServiceSSL {
		err1 := usbMuxClient.deviceConnection.EnableSessionSSL(pairRecord)
		if err1 != nil {
			return err1
		}
	}
	return nil
}
