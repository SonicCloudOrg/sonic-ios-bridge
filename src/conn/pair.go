package conn

import (
	"bytes"
	"fmt"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/tool"
	"howett.net/plist"
)

type PairRecord struct {
	HostID            string
	SystemBUID        string
	HostCertificate   []byte
	HostPrivateKey    []byte
	DeviceCertificate []byte
	EscrowBag         []byte
	WiFiMACAddress    string
	RootCertificate   []byte
	RootPrivateKey    []byte
}

type PairRecordData struct {
	PairRecordData []byte
}

type ReadPairMessage struct {
	BundleID            string
	ClientVersionString string
	MessageType         string
	ProgName            string
	LibUSBMuxVersion    uint32 `plist:"kLibUSBMuxVersion"`
	PairRecordID        string
}

func NewReadPairMessage(udId string) ReadPairMessage {
	data := ReadPairMessage{
		MessageType:         "ReadPairRecord",
		PairRecordID:        udId,
		BundleID:            BundleId,
		ClientVersionString: ClientVersion,
		ProgName:            ProgramName,
		LibUSBMuxVersion:    3,
	}
	return data
}

func (usbMuxClient *UsbMuxClient) ReadPair(udId string) (PairRecord, error) {
	err := usbMuxClient.Send(NewReadPairMessage(udId))
	if err != nil {
		return PairRecord{}, tool.NewErrorPrint(tool.ErrConnect, "usbMux", err)
	}
	resp, err := usbMuxClient.ReadMessage()
	if err != nil {
		return PairRecord{}, tool.NewErrorPrint(tool.ErrReadingMsg, "pair", err)
	}
	pairRecordData, _ := pairRecordDataForBytes(resp.Payload)
	return pairRecordForBytes(pairRecordData.PairRecordData), nil
}

func pairRecordDataForBytes(plistBytes []byte) (PairRecordData, error) {
	decoder := plist.NewDecoder(bytes.NewReader(plistBytes))
	var data PairRecordData
	decoder.Decode(&data)
	if data.PairRecordData == nil {
		resp := usbMuxRespForBytes(plistBytes)
		return data, fmt.Errorf("device not pair : %d", resp.Number)
	}
	return data, nil
}

func pairRecordForBytes(plistBytes []byte) PairRecord {
	decoder := plist.NewDecoder(bytes.NewReader(plistBytes))
	var data PairRecord
	decoder.Decode(&data)
	return data
}
