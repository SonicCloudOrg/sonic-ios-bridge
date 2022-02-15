package conn

import (
	"bytes"
	"howett.net/plist"
)

type StartSessionRequest struct {
	Label           string
	ProtocolVersion string
	Request         string `plist:"Request"`
	HostID          string
	SystemBUID      string
}

type StartSessionResp struct {
	EnableSessionSSL bool
	Request          string
	SessionID        string
}

type StopSessionRequest struct {
	Request   string `plist:"Request"`
	SessionID string
	Label     string
}

type StopSessionResp struct {
	Request string
}

func NewStartSessionRequest(hostID string, systemBUID string) StartSessionRequest {
	return StartSessionRequest{
		Request:         "StartSession",
		Label:           BundleId,
		ProtocolVersion: "2",
		HostID:          hostID,
		SystemBUID:      systemBUID,
	}
}

func NewStopSessionRequest(sessionID string) StopSessionRequest {
	data := StopSessionRequest{
		Request:   "StopSession",
		SessionID: sessionID,
		Label:     BundleId,
	}
	return data
}

func startSessionRespForBytes(plistBytes []byte) StartSessionResp {
	decoder := plist.NewDecoder(bytes.NewReader(plistBytes))
	var data StartSessionResp
	_ = decoder.Decode(&data)
	return data
}

func (lockDownConn *LockDownConnection) StartSession(pairRecord PairRecord) (StartSessionResp, error) {
	err := lockDownConn.Send(NewStartSessionRequest(pairRecord.HostID, pairRecord.SystemBUID))
	if err != nil {
		return StartSessionResp{}, err
	}
	resp, err := lockDownConn.ReadMessage()
	if err != nil {
		return StartSessionResp{}, err
	}
	resp1 := startSessionRespForBytes(resp)
	lockDownConn.sessionID = resp1.SessionID
	if resp1.EnableSessionSSL {
		err = lockDownConn.deviceConnection.EnableSessionSSL(pairRecord)
		if err != nil {
			return StartSessionResp{}, err
		}
	}
	return resp1, nil
}

func (lockDownConn *LockDownConnection) StopSession() {
	if lockDownConn.sessionID == "" {
		return
	}
	lockDownConn.Send(NewStopSessionRequest(lockDownConn.sessionID))
	lockDownConn.ReadMessage()
}
