package conn

import (
	"bytes"
	"howett.net/plist"
)

type UsbMuxResponse struct {
	MessageType string
	Number      uint32
}

func UsbMuxRespForBytes(plistBytes []byte) UsbMuxResponse {
	decoder := plist.NewDecoder(bytes.NewReader(plistBytes))
	var usbMuxResponse UsbMuxResponse
	_ = decoder.Decode(&usbMuxResponse)
	return usbMuxResponse
}

// IsSuccess https://github.com/mogaleaf/java-usbmuxd
func (usbMuxResponse UsbMuxResponse) IsSuccess() bool {
	return usbMuxResponse.Number == 0
}
