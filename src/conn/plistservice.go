package conn

import (
	"fmt"
	plist "howett.net/plist"
)

const (
	MessageTypeDeviceList     = "ListDevice"
	MessageTypeReadPairRecord = "ReadPairRecord"
)

func TransToPlistBytes(data interface{}) []byte {
	bytes, err := plist.Marshal(data, plist.XMLFormat)
	if err != nil {
		fmt.Errorf("fail converting to plist %v :%w", data, err)
	}
	return bytes
}
