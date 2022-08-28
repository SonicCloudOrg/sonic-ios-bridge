package webinspector

import (
	"fmt"
	giDevice "github.com/electricbubble/gidevice"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"testing"
)

func TestCDPWebSocket(t *testing.T) {
	r := gin.Default()
	giDevice.SetDebug(true, true)
	SetProtocolDebug(true)
	r.GET("/", pagesHandle)
	r.GET("/json", pagesHandle)
	r.GET("/json/list", pagesHandle)
	r.GET("/devtools/page/:id", pageDebugHandle)

	r.Run(fmt.Sprintf("127.0.0.1:%d", port))
}

func TestGetWSData(t *testing.T) {
	// ios web inspector proc
	dialer := websocket.Dialer{}
	connect, _, err := dialer.Dial("ws://127.0.0.1:9222/devtools/page/1", nil)
	if nil != err {
		log.Println(err)
		return
	}
	//离开作用域关闭连接，go 的常规操作
	defer connect.Close()

	//启动数据读取循环，读取客户端发送来的数据
	for {
		//从 websocket 中读取数据
		//messageType 消息类型，websocket 标准
		//messageData 消息数据
		messageType, messageData, err := connect.ReadMessage()
		if nil != err {
			log.Println(err)
			break
		}
		switch messageType {
		case websocket.TextMessage: //文本数据
			// {"method":"Target.targetCreated","params":{"targetInfo":{"targetId":"page-31","type":"page"}}}
			fmt.Println(string(messageData))
			fmt.Println("text data")
		case websocket.BinaryMessage: //二进制数据
			fmt.Println("binary data")
			fmt.Println(messageData)
		case websocket.CloseMessage: //关闭
			fmt.Println("end close")
		case websocket.PingMessage: //Ping
			fmt.Println("ping service")
		case websocket.PongMessage: //Pong
			fmt.Println("pong service")
		default:

		}
	}
}
