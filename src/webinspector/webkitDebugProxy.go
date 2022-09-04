package webinspector

import (
	"context"
	"fmt"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	giDevice "github.com/electricbubble/gidevice"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var webDebug *WebkitDebugService
var localPort = 9222

func InitWebInspectorServer(udid string, port int, isDebug bool) context.CancelFunc {
	var err error
	var cannel context.CancelFunc
	if webDebug == nil {
		// 优化初始化过程
		ctx := context.Background()
		device := util.GetDeviceByUdId(udid)
		webDebug = NewWebkitDebugService(&device, ctx)
		cannel, err = webDebug.ConnectInspector()
		if err != nil {
			log.Fatal(err)
		}
	}
	localPort = port
	if isDebug {
		SetProtocolDebug(true)
		giDevice.SetDebug(true, true)
	}
	return cannel
}

func PagesHandle(c *gin.Context) {

	pages, err := webDebug.GetOpenPages(localPort)
	if err != nil {
		c.JSONP(http.StatusNotExtended, err)
	}
	c.JSONP(http.StatusOK, pages)
}

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options

func PageDebugHandle(c *gin.Context) {
	id := c.Param("id")

	application, page, err := webDebug.FindPagesByID(id)
	if application == nil || page == nil {
		c.Error(fmt.Errorf(fmt.Sprintf("not find page to id:%s", id)))
		log.Println(fmt.Errorf(fmt.Sprintf("not find page to id:%s", id)))
		return
	}
	conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Print("Error during connection upgradation:", err)
		return
	}
	defer conn.Close()

	go func() {
		err = webDebug.StartCDP(application.ApplicationID, page.PageID, conn)
		if err != nil {
			log.Fatal(err)
		}
	}()
	//// 确保初始化完成
	err = webDebug.ReceiveProtocolData()
	if err != nil {
		fmt.Println(err)
	}
	go func() {
		for {
			err = webDebug.ReceiveProtocolData()
			if err != nil {
				return
			}
		}
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error during message reading:", err)
			break
		}
		if message != nil {
			if len(message) == 0 {
				continue
			}
			webDebug.SendProtocolCommand(application.ApplicationID, page.PageID, message)
		}
	}
}
