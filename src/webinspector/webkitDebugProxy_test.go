package webinspector

import (
	"encoding/json"
	"fmt"
	"github.com/SonicCloudOrg/sonic-ios-webkit-adapter/entity"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"testing"
)

func TestCDPWebSocket(t *testing.T) {
	r := gin.Default()
	//giDevice.SetDebug(true, true)
	SetProtocolDebug(true)
	r.GET("/", PagesHandle)
	r.GET("/json", PagesHandle)
	r.GET("/json/list", PagesHandle)
	r.GET("/devtools/page/:id", PageDebugHandle)

	r.Run(fmt.Sprintf("127.0.0.1:%d", localPort))
}

func TestGetWSData(t *testing.T) {
	// ios web inspector proc
	dialer := websocket.Dialer{}
	connect, _, err := dialer.Dial("ws://127.0.0.1:9222/devtools/page/1", nil)
	if nil != err {
		log.Println(err)
		return
	}
	defer connect.Close()

	for {
		messageType, messageData, err := connect.ReadMessage()
		if nil != err {
			log.Println(err)
			break
		}
		switch messageType {
		case websocket.TextMessage:
			// {"method":"Target.targetCreated","params":{"targetInfo":{"targetId":"page-31","type":"page"}}}
			fmt.Println(string(messageData))
			fmt.Println("text data")
		case websocket.BinaryMessage:
			fmt.Println("binary data")
			fmt.Println(messageData)
		case websocket.CloseMessage:
			fmt.Println("end close")
		case websocket.PingMessage:
			fmt.Println("ping service")
		case websocket.PongMessage:
			fmt.Println("pong service")
		default:

		}
	}
}

func TestJsonToStruct(t *testing.T) {
	msg := "{\"id\":15,\"method\":\"Log.enable\",\"params\":{}}"
	protocolMessage := &entity.TargetProtocol{}
	err := json.Unmarshal([]byte(msg), protocolMessage)
	if err != nil {
		log.Panic(err)
	}
}
