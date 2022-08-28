package webinspector

import (
	"fmt"
	"github.com/SonicCloudOrg/sonic-ios-bridge/src/util"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var webDebug *WebkitDebugService
var port = 9222

// todo add lock
func pagesHandle(c *gin.Context) {
	if webDebug == nil {
		// 优化初始化过程
		device := util.GetDeviceByUdId("")
		webDebug = NewWebkitDebugService(&device)
		err := webDebug.ConnectInspector()
		if err != nil {
			panic(err)
		}
	}
	pages, err := webDebug.GetOpenPages(port)
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

func pageDebugHandle(c *gin.Context) {
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
		err = webDebug.StartCDP(application.ApplicationID, page.PageID)
		if err != nil {
			log.Fatal(err)
		}
	}()
	//// 确保初始化完成
	webDebug.ReceiveProtocolData(conn)
	go func() {
		for {
			webDebug.ReceiveProtocolData(conn)
		}
	}()

	for {
		//todo target SendF
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error during message reading:", err)
			break
		}
		if message != nil {
			if len(message) == 0 {
				continue
			}
			go webDebug.SendProtocolCommand(application.ApplicationID, page.PageID, message)
		}
	}
}
