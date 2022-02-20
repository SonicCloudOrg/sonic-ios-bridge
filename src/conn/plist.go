package conn

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	plist "howett.net/plist"
	"io"
)

type PlistCodec struct {
}

func NewPlistCodec() PlistCodec {
	return PlistCodec{}
}

func transToPlistBytes(data interface{}) []byte {
	bytes, err := plist.Marshal(data, plist.XMLFormat)
	if err != nil {
		fmt.Errorf("fail transfer to plist %v :%w", data, err)
	}
	return bytes
}

func transToPlistString(data interface{}) string {
	bytes, _ := plist.Marshal(data, plist.XMLFormat)
	return string(bytes)
}

func parsePlist(data []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	_, err := plist.Unmarshal(data, &result)
	return result, err
}

func (plistCodec PlistCodec) Encode(message interface{}) ([]byte, error) {
	stringContent := transToPlistString(message)
	buf := new(bytes.Buffer)
	length := len(stringContent)
	messageLength := uint32(length)

	err := binary.Write(buf, binary.BigEndian, messageLength)
	if err != nil {
		return nil, err
	}
	buf.Write([]byte(stringContent))
	return buf.Bytes(), nil
}

func (plistCodec PlistCodec) Decode(r io.Reader) ([]byte, error) {
	if r == nil {
		return nil, errors.New("reader equals nil")
	}
	buf := make([]byte, 4)
	err := binary.Read(r, binary.BigEndian, buf)
	if err != nil {
		return nil, err
	}
	size := binary.BigEndian.Uint32(buf)
	payLoadBytes := make([]byte, size)
	_, err = io.ReadFull(r, payLoadBytes)
	if err != nil {
		return nil, err
	}
	return payLoadBytes, nil
}
