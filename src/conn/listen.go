package conn

import (
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/tool"
)

type ListenMessage struct {
	MessageType         string
	ProgName            string
	ClientVersionString string
	ConnType            int
	kLibUSBMuxVersion   int
}

func NewListen() ListenMessage {
	msg := ListenMessage{
		MessageType:         "Listen",
		ProgName:            ProgramName,
		ClientVersionString: ClientVersion,
		ConnType:            1,
		kLibUSBMuxVersion:   3,
	}
	return msg
}

func (usbMuxClient *UsbMuxClient) Listen() (func() (iDevice, error), error) {
	msg := NewListen()
	err := usbMuxClient.Send(msg)
	if err != nil {
		return nil, err
	}
	resp, err := usbMuxClient.ReadMessage()
	if err != nil {
		return nil, err
	}
	if !usbMuxRespForBytes(resp.Payload).IsSuccess() {
		return nil, tool.NewErrorPrint(tool.ErrSendCommand, "listen", err)
	}
	return func() (iDevice, error) {
		usb, err := usbMuxClient.ReadMessage()
		if err != nil {
			return iDevice{}, err
		}
		return deviceForBytes(usb.Payload), nil
	}, nil

}
