package conn

import "fmt"

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
	if !UsbMuxRespForBytes(resp.Payload).IsSuccess() {
		return nil, fmt.Errorf("failed to send listen command : %w", err)
	}
	return func() (iDevice, error) {
		usb, err := usbMuxClient.ReadMessage()
		if err != nil {
			return iDevice{}, err
		}
		return DeviceForBytes(usb.Payload), nil
	}, nil

}
