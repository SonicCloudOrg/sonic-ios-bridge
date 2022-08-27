package webinspector

import (
	"encoding/json"
	"fmt"
	"log"
)

type SessionProtocol struct {
	appID       string
	pageID      int
	webDebug    *WebkitDebugServer
	inputEvent  chan []byte
	outputEvent chan []byte
}

var isProtocol = false

func (p *SessionProtocol) SendMessage(message []byte) {
	if isProtocol {
		log.Println(fmt.Sprintf("protocol send command :%s", string(message)))
	}
	if arr, err := json.Marshal(message); err != nil {
		log.Fatal(err)
	} else {
		p.inputEvent <- arr
	}
}

func (p *SessionProtocol) loopInput() {
	go func() {
		for {
			select {
			case data, ok := <-p.inputEvent:
				if ok {
					go p.webDebug.SendCommand(p.appID, p.pageID, data)
				}
			}
		}
	}()
}

func (p *SessionProtocol) RecvMessage() []byte {
	return p.webDebug.RecvProtocolData()
}
